# my-texas-42-backend

Go HTTP + WebSocket server for MyTexas42, a multiplayer Texas 42 domino game. Self-hosted via Docker Compose (see `../mytexas42-compose`). Replaces the legacy AWS-Lambda+DynamoDB backend in `../my-texas-42-react-app/packages/functions`.

**Companion frontend:** `../my-texas-42-frontend`.

**Status:** Scaffolding, auth, lobby, chat, and friends are done. **The Texas 42 game engine (bid/call/play validation, trick-taking, scoring, round transitions) is stubbed out** — `models.GameState.ProcessMove()` accepts anything. See `MIGRATION_TODO.md` for the full porting checklist.

## Stack

- **Go 1.21**, single binary
- **Gin** for HTTP routing + middleware (`main.go`)
- **gorilla/websocket** for the `/ws` realtime channel
- **PostgreSQL** via `lib/pq` (no ORM — bespoke reflection helper in `services/database.go`, raw SQL builders in `sql_scripts/`)
- **AWS Cognito** still used for identity (User Pool from the legacy stack); `services/cognito.go` calls Cognito directly per request
- Logs persist to a `Logs` table via `logger/log.go`

## Layout

```
my-texas-42-backend/
├── main.go                 # Routes, CORS, server bootstrap
├── admin/                  # GET /app stats (admin-only)
├── auth/                   # Authenticate middleware (HTTP + WS branches), admin guard
├── friends/                # POST/DELETE friends, accept request
├── games/                  # GET/POST /games, in-memory GameManager registry
├── logger/                 # log.go — writes INFO/WARN/ERROR/CRITICAL into Logs table
├── models/
│   ├── object-models.go    # GameState, GlobalGameState, PlayerGameState, DominoName
│   ├── game-extensions.go  # NewGame, AssignDominoes, ProcessMove (STUB), rule constants
│   ├── game-validations.go # validateTurn (stub-level)
│   ├── api-models.go       # Request/response shapes
│   └── decode-api-model.go # Generic JSON request decoder
├── request-util/           # Pulls username/sub/email out of gin.Context
├── services/
│   ├── cognito.go          # Login, SignUp, Confirm, ChangePassword, AuthenticateRequest
│   └── database.go         # Query[T] reflection helper, Execute
├── sockets/                # WS connect, message dispatch, chat, switch-teams, refresh, process-move
├── sql_scripts/            # Hand-rolled SQL builders (sanitized via util.Sanitize, NOT parameterized)
├── system/initialize.go    # Env var loading; panics if required vars missing
├── users/                  # Signup, ConfirmSignup, Login, GetCurrentUser, GetUserProfile, ChangePassword, ChangeDisplayName
└── util/                   # Sanitize + IsXValid validators
```

## Routes

```
GET    /health
GET    /app                                  (auth + admin)
POST   /users
PUT    /users/confirm
POST   /users/resend-confirmation
POST   /users/login
GET    /users/current                        (auth)
PUT    /users/change-password                (auth)
PUT    /users/change-display-name            (auth)
GET    /users/:username                      (auth)
POST   /friends/:username                    (auth)
POST   /friends/:username/accept             (auth)
DELETE /friends/:username                    (auth)
GET    /games                                (auth)
POST   /games                                (auth)
GET    /ws                                   (auth, upgrades to WebSocket)
```

WebSocket message envelopes (snake_case both directions):

- **Client → server:** `{"action": string, "data": string}` where `action ∈ {send_chat_message, play_turn, refresh_player_game_state, switch_teams}`. For `play_turn`, `data` is a slash-delimited string: `"bid/31"`, `"call/Nil"`, `"play/6:4"`.
- **Server → client:** `{"message_type": "chat"|"game-update"|"game-error", "message": string, "username": string, "game_data": PlayerGameState}`.

The frontend depends on the snake_case keys above. If you add fields to `GameState`/`PlayerGameState`, **set explicit `json:"..."` tags in snake_case**.

## Persistence

- **PostgreSQL** schema is defined in `sql_scripts/table-schema.go` (and shadowed by `../mytexas42-compose/postgres-build/mytexas42-schema.sql` — the Go file wins because the schema is created from code).
- Tables: `Users`, `Friends`, `FriendRequests`, `UserStats`, `MatchArchive`, `RoundArchive`, `ChatMessageArchive`, `Logs`.
- Active games live **in memory only** in `games.GameManager.games` keyed by invite code. **A restart drops every active match.**
- `services/database.go::Query[T]` opens a fresh connection per call (no pooling), executes raw SQL, scans results positionally into struct fields.
- All SQL goes through `util.Sanitize()` (escapes `'`, strips `;` and `\n`). This is **not** equivalent to parameterized queries — be deliberate; treat new query builders as a security review surface.

## Auth

- HTTP: `Authorization: <Cognito access token>` → middleware calls `services.AuthenticateRequest()` → Cognito `GetUser` → DB lookup → context (`user`, `sub`, `emailVerified`, `email`).
- WebSocket: `?username=&sub=` query params; middleware skips Cognito and resolves the user by `sub` from the DB. **Trust boundary is weaker than HTTP** — anyone with a known `sub` can impersonate. Don't widen this.
- Cognito user pools are per-environment; configured by env vars (`STAGING_USER_POOL_*`, `PRODUCTION_USER_POOL_*`).

## Run / build

```bash
# Local dev (requires Postgres reachable + env vars from ../mytexas42-compose/.env)
go run .

# Container build (multi-stage, scratch runtime)
docker build -t mytexas42-backend .

# Typical workflow: bring up everything together
cd ../mytexas42-compose && docker compose up -d backend-staging
```

Required env vars (see `system/initialize.go` for the exhaustive list): `ENVIRONMENT`, `PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `POSTGRES_*_HOST_NAME`, `STAGING_USER_POOL_NAME`/`_APP_KEY`, `PRODUCTION_USER_POOL_NAME`/`_APP_KEY`. The process **panics** if any are missing.

## Conventions

- Handler functions are gin handlers: `func(c *gin.Context)`. Read the user via `request-util.GetUser(c)`.
- Errors back to clients: `c.JSON(status, gin.H{"error": "msg"})`. No typed error hierarchy.
- Logs: `logger.Log(level, "service-name", msg)` — these become DB rows. Don't put high-frequency events through this.
- New SQL: add a function in `sql_scripts/` that builds a string and call via `services.Query[T]` or `services.Execute`. Always sanitize *every* string interpolated.
- New WebSocket actions: add a case in `sockets/handle-incoming-messages.go` (the dispatcher) and a handler file alongside `handle-chat-message.go`.

## Known gaps / hazards

The biggest one: **`models.GameState.ProcessMove()` is a stub**. It parses the move string and returns nil. There is no bid validation, no trick winner logic, no scoring, no round/match transitions. The legacy implementation in `../my-texas-42-react-app/packages/functions/src/utils/game-utils.ts` is the spec. See `MIGRATION_TODO.md`.

Other things to know:

- No tests (`*_test.go` count: 0).
- Game state isn't persisted between restarts.
- `UserStats` rows are read but never updated.
- `RoundArchive` is never written to; `MatchArchive` is created at game start but never finalized with winners/scores.
- No connection pool, no DB transactions.
- No WebSocket heartbeat / dead-connection detection.
- No rate limiting anywhere.
- `getMove` parses moves but `ProcessMove` discards the result without enforcing rules.
