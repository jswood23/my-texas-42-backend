# MyTexas42 Backend — Migration Gap Spec

Hand-off doc for porting the legacy AWS-Lambda backend (`../my-texas-42-react-app/packages/functions`) to this Go service. Pair this with `CLAUDE.md` (architecture overview) and the legacy reference: `../my-texas-42-react-app/packages/functions/src/utils/game-utils.ts` is the **canonical rules implementation** — port its semantics, not its style.

Work items are ordered roughly by criticality. P0 items block real games from being playable; P1 items block durability + correctness; P2 items are quality.

---

## P0-A. Implement the Texas 42 game engine

**Where it lives today:** `models/game-extensions.go::GameState.ProcessMove(username, moveStr)` parses the move string into `(MoveType, ActualMove)` and returns `nil`. There is no validation or state mutation. `models/game-validations.go::validateTurn` only checks "is it your turn".

**What needs to happen:** every accepted move must update game state correctly and reject invalid moves with a `game-error` to the sender (other clients should not see an error).

### Acceptance criteria

The behavior must match `utils/game-utils.ts::checkValidity`, `processBids`, `setRoundRules`, `playDomino`, `getWinningPlayerOfTrick`, `getTrickScore`, `processEndOfTrick`, `processRoundWinner`, `startNextRound`, `skipPlayerTurnIfNil`, `getPlayerPosition` from the legacy code.

### Suggested file plan

```
models/
├── moves.go        // MoveType, ActualMove, parseMove (already partly in game-extensions.go)
├── bidding.go      // validateBid, processBids
├── calling.go      // validateCall, applyCall  (sets RoundRules.Trump + .Variant)
├── play.go         // validatePlay, applyPlay  (suit-following + trump priority + variant rules)
├── trick.go        // determineTrickWinner (per-variant table), trickScore
├── scoring.go      // endOfTrick, endOfRound, awardMarks, isMatchOver
└── round.go        // startNextRound (rotate starting bidder, redeal)
```

Keep `GlobalGameState` immutable inside validators (return a new state or an error); apply mutations only after validation succeeds. The legacy code mutates in place — don't replicate that.

### Bid validation rules (from legacy `checkValidity`)

- **Pass** is `bid == 0`. Always allowed.
- **Numeric bid** must be `>= 30 && <= 41`, strictly greater than the current high.
- **Mark bid** must be a multiple of 42 (`42, 84, 126, 168`). Can only jump **one** mark at a time above the previous high (`prev_high == 41 → 42 OK, 84 not OK`; `prev_high == 42 → 84 OK, 126 not OK`).
- **Forced 31 bid** rule (when `Rules` contains `"Forced 31 bid"`): if the first three players pass, the dealer's minimum bid is 31, not 30, and pass is disallowed.
- **Splash** (mark bid + variant chosen at call time) requires the bidder to hold ≥3 doubles.
- **Plunge** requires ≥4 doubles and is a 2-mark bid that rotates to the partner to pick trump.

### Call validation rules (`setRoundRules` in legacy)

After bids, the bidding team's bidder calls trump + variant. Variants:

| Variant | Trump | Notes |
|---|---|---|
| Numeric (0–6) | the suit number | doubles belong to their suit by default |
| `"Follow-Me"` | trump = led suit each trick | doubles can be high or low |
| `"Doubles-Trump"` | doubles **are** the trump suit | |
| `"Nil"` | no trump | bidder's partner sits out (3-player tricks); Nil bidder must take **zero** tricks |
| Nil sub-variants | `"Doubles-high"` / `"Doubles-low"` / `"Doubles-own-suit"` | how doubles compare against non-doubles |
| `"Splash"` | bidder declares trump from their hand | 2-mark, ≥3 doubles |
| `"Plunge"` | partner declares trump | 2-mark, ≥4 doubles, partner plays as bidder |
| `"Sevens"` | no trump | each trick won by domino whose pip-sum is closest to 7 |

Reject any call that the rule set / bid type doesn't permit (e.g., no `"Splash"` if Splash isn't in `Rules`).

### Play validation rules

Inputs: current trick, player's hand, `RoundRules`, `Rules`. Reject if:

1. Player doesn't hold the domino.
2. There is a led suit and the player has a domino of that suit but plays a different non-trump non-led-suit one. (Trump always overrides led suit, except in Follow-Me where the led suit *is* the trump unless overridden by the variant; doubles' "suit" depends on variant.)
3. Sevens variant: player must play the domino in their hand whose pip-sum is closest to 7 if such a domino exists and matches the closeness constraint defined in legacy code.
4. Doubles-Own-Suit Nil variant: doubles are treated as their own suit; led-suit rules differ.

The legacy `checkValidity()` is the source of truth for the suit-matching matrix.

### Trick winner determination

Implement as a function table keyed by `(trump, variant)`:

| Case | Winner |
|---|---|
| Numeric trump played | highest trump (doubles count as their own suit's pip-pair) |
| No trump played | highest in the led suit |
| Doubles-Trump | highest double, else highest in led suit |
| Follow-Me | highest in led suit (doubles configurable per sub-rule) |
| Nil + Doubles-high | doubles beat anything in led suit |
| Nil + Doubles-low | doubles lose to any non-double in led suit |
| Nil + Doubles-own-suit | doubles only win their suit's tricks |
| Sevens | domino with pip-sum closest to 7 |

For the Nil variants the bidder's *partner* never plays (turn skipping), so tricks have 3 dominoes — handle empty seat gracefully.

### Trick scoring

A domino "counts" if its pip-sum is divisible by 5: `5-0, 4-1, 3-2, 6-4, 5-5, 6-3` worth their pip-sum. Total point dominoes: 35 across the deck + 1 per trick taken = 42 points per round. (Match the legacy `getTrickScore`.)

### Round-end detection

After each trick, check:

- **Bidding team meets bid:** their cumulative round score ≥ `RoundRules.Bid` (or for mark bids: they took every trick, depending on variant).
- **Non-bidding team breaks the bid:** their cumulative round score > `42 - RoundRules.Bid` (round score is always 0–42).
- **Nil:** Nil bidder takes any trick → bidding team loses immediately.
- **Sevens:** Round runs to 7 tricks; team with more tricks wins (handle ties per legacy).

On round end, award marks (`ceil(bid / 42)`) to the winning team. Then call `startNextRound`: rotate starting bidder by one player, reset round score, redeal via `AssignDominoes`.

### Match end

Legacy has **no match-end condition** (games loop forever). The new server should add one — typical Texas 42 plays to **7 marks**. Configurable via a `Rules` entry, default 7. On match end, finalize `MatchArchive` (set winners, total scores) and write out per-round `RoundArchive` rows.

---

## P0-B. JSON tags on game-state structs

**Why:** the frontend expects snake_case keys (`match_name`, `current_player_turn`, `current_round_rules`, etc.) on every WebSocket `game_data` payload. Go's default marshaling produces PascalCase.

**Action:** add `json:"snake_case"` tags to every exported field on `GameState`, `GlobalGameState`, `PlayerGameState`, and `RoundRules` in `models/object-models.go` (and any nested types). The frontend `GameState` type in `../my-texas-42-frontend/src/types/index.d.ts` is the contract.

**Field mapping (non-exhaustive):**

```
MatchName                 → match_name
MatchInviteCode           → match_invite_code
MatchPrivacy              → match_privacy
Rules                     → rules
Team1UserNames            → team_1
Team2UserNames            → team_2
Team1Connected            → team_1_connected
Team2Connected            → team_2_connected
CurrentRound              → current_round
CurrentStartingBidder     → current_starting_bidder
CurrentStartingPlayer     → current_starting_player
IsInBidding               → current_is_bidding   ⚠ name change
CurrentPlayerTurn         → current_player_turn
RoundRules                → current_round_rules  ⚠ name change + JSON-string-encoded
Team1RoundScore           → current_team_1_round_score
Team2RoundScore           → current_team_2_round_score
CurrentTeam1TotalScore    → current_team_1_total_score
CurrentTeam2TotalScore    → current_team_2_total_score
RoundHistory              → current_round_history
TotalRoundHistory         → total_round_history
PlayerDominoes            → player_dominoes
HasStarted                → has_started
```

`current_round_rules` is **a JSON string** in the legacy contract, and the new frontend already calls `JSON.parse` on it. The simplest port is to marshal `RoundRules` to a string before assembling the outbound message in `sockets/message-players-in-game.go`. (Alternative: change the contract to a nested object and update both sides — easier long-term, but more work right now.)

Verify by capturing one `game-update` message in the browser DevTools Network panel and diffing against the legacy production payload.

---

## P0-C. Move-history wire format

The legacy frontend reads `current_round_history` as `string[]` where each entry is `username\moveType\move` (backslash-delimited). The new frontend in `pages/play/in-game/game-window/utils/get-game-information.ts` does the same parsing. Make sure `ProcessMove` appends in this exact format. **Don't switch delimiters** — the frontend split logic is already wired.

The new in-bound action format is `moveType/move` (slash-delimited from the client). On the server, prepend `username\` and switch the slash to a backslash before appending to history.

---

## P0-D. Missing handlers / response shapes the frontend already calls

Cross-check against `../my-texas-42-frontend/src/utils/api-utils.ts` callers.

| Endpoint | Status | Notes |
|---|---|---|
| `POST /users/login` | exists | confirm response shape `{message, token}` matches `LoginResponseAPIModel` |
| `POST /users/signup` | route is `POST /users` today | **rename** route or update frontend |
| `POST /users/{username}/confirm` | route is `PUT /users/confirm` (no path param) | unify |
| `POST /users/{username}/resend-confirmation` | route is `POST /users/resend-confirmation` | unify |
| `PUT /users/profile` | route is `PUT /users/change-display-name` | unify |
| `PUT /users/password` | route is `PUT /users/change-password` | unify |
| `GET /games` | exists | confirm response includes `in_game` (current game's invite code if any), `public_games`, `private_games` |
| `POST /games` | exists | confirm response is `{invite_code, team_number}` |
| `POST /friends/{username}` | exists (AddFriend) | sends or auto-accepts mutual |
| `POST /friends/{username}/accept` | exists | dedicated accept endpoint |
| `DELETE /friends/{username}` | exists (RemoveFriendOrRequest) | covers cancel, reject, remove — confirm all three branches |

**Action:** decide whether to rename Go routes to match the frontend, or update the frontend's `api-utils.ts` callers. Recommend renaming the backend routes — the frontend names are cleaner and the backend has fewer call sites.

---

## P1-A. Persist active game state

Active games live in `games.GameManager.games` (in-memory). A backend restart drops every match. Players reconnecting will get a "lobby not found" error.

**Spec:**

- After every state mutation in `ProcessMove` (and `connect`, `disconnect`, `switch-teams`), upsert the `GlobalGameState` to a new `ActiveMatch` table.
- On startup, `games.GameManager` rehydrates from `ActiveMatch`.
- On match end, delete the `ActiveMatch` row (the snapshot is preserved via `MatchArchive` + `RoundArchive`).
- Wrap the upsert and any related table writes in a `database/sql` transaction.

**Schema:**

```sql
CREATE TABLE ActiveMatch (
  InviteCode VARCHAR(6) PRIMARY KEY,
  MatchID INT NOT NULL REFERENCES MatchArchive(MatchID),
  StateJSON JSONB NOT NULL,
  UpdatedAt TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON ActiveMatch (MatchID);
```

Persist as JSONB; rehydration is `json.Unmarshal` into `GlobalGameState`. Keep the in-memory map as the read path, the DB as the durability layer (write-through).

---

## P1-B. Update UserStats and finalize MatchArchive / RoundArchive

Currently `UserStats` rows are read on profile fetch but never updated. `MatchArchive` is created at game start but never marked complete. `RoundArchive` is never written.

**Spec:**

At each round end (`endOfRound` in the new engine):

1. Insert `RoundArchive` row: `MatchID`, `RoundRules` (bit-packed or JSON-encoded), `Team1Score`, `Team2Score`, `RoundActivity` (a serialized round history string).
2. Increment `RoundsPlayed` for all four players. Increment `RoundsWon` for the two on the winning team.
3. Bidder accumulates `TotalPointsAsBidder` += team round score; `TotalRoundsAsBidder` += 1; `TimesWinningBidTotal`/`TimesCallingSuit`/`TimesCallingNil`/etc. by variant.
4. Bidder's partner accumulates `TotalPointsAsSupport` / `TotalRoundsAsSupport`.
5. Defending team accumulates `TotalPointsAsCounter` / `TotalRoundsAsCounter`.

At match end:

1. Update `MatchArchive` with `WinningTeam`, `Team1Marks`, `Team2Marks`, `TotalRounds`.
2. Increment `GamesPlayed` for all four; `GamesWon` for winners.

Wrap all per-round and per-match writes in transactions.

---

## P1-C. WebSocket disconnect during play

`sockets/disconnect.go` removes the connection but doesn't change game state. If a player drops mid-round, the game just stalls on their turn.

**Spec:** when a player disconnects:

- Mark `Team1Connected[i]` / `Team2Connected[i]` false; broadcast `game-update`.
- If the disconnect happens on their turn, start a 60-second reconnect timer (in-memory).
- On reconnect within the window: resume.
- On timeout: forfeit the round to the other team (award them remaining marks needed to break the bid). Match continues.
- If both players on a team time out: forfeit the match.

The legacy backend also doesn't handle this — design fresh.

---

## P1-D. Replace string-formatted SQL with parameterized queries

`sql_scripts/*.go` builds SQL strings via `fmt.Sprintf` and runs `util.Sanitize()` on inputs. That's not a real defense against SQL injection — it's an escape function for single quotes.

**Spec:** switch to `db.Query(query, args...)` / `db.Exec(query, args...)` with `$1, $2, ...` placeholders. Drop `util.Sanitize` calls from query builders (keep it only for chat content if you still want to filter `;`, `'`, `\n` from rendered messages — that's a UX choice, not a security one).

This is a security-prioritized refactor; do it before exposing any new SQL endpoints to user input.

---

## P1-E. WebSocket auth strengthening

The `/ws` middleware skips Cognito and resolves the user by `?sub=...` from the DB. Anyone who knows a user's `sub` can connect as them.

**Spec:** require the same `Authorization: <Cognito access token>` header on the WebSocket upgrade request as on HTTP routes (or accept it as a `?token=` query string and validate via Cognito `GetUser` once at connect time). Reject if the resolved username doesn't match the `?username=` claim.

## P1-F. DB connection pooling

`services/database.go::Query[T]` opens a connection per call. Move to a long-lived `*sql.DB` (connection pool is built into `database/sql`) initialized in `system.Initialize()` and shared via context or a package-level var.

---

## P2. Quality items

- **Tests.** No `*_test.go` files exist. The game engine is the highest-leverage place — a table-driven test suite that runs known-good move sequences against the engine and asserts state will catch most regressions. Use the legacy tests as inspiration only (there aren't any), but generate fixtures from running the legacy engine against curated scenarios.
- **Heartbeat / dead-connection detection.** gorilla/websocket supports ping/pong; current code has none. A 30s ping with 60s read deadline is standard.
- **Rate limiting.** None on REST or WS. Minimum: per-IP signup limit, per-user chat-message rate limit (say, 5 / 5s).
- **Structured logging.** `logger/log.go` writes string rows to a DB table. Use a real logger (`slog`) and ship to stdout — the DB log table fills up fast.
- **Migrations.** Schema is created from Go code in `sql_scripts/table-schema.go`. Adopt `golang-migrate` or `goose` so we can roll forward without coordinated deploys.
- **Connection registry race conditions.** `sockets.ConnectionManager` uses a single `sync.Mutex`. For higher concurrency, switch to `sync.Map` or shard by username hash.
- **Game listing pagination.** `GET /games` returns everything. Add `?limit=&offset=`.

---

## Reference: legacy spec files

The fastest way to clarify any rules-question is to read the legacy implementation:

- Engine: `../my-texas-42-react-app/packages/functions/src/utils/game-utils.ts`
- Lobby ops: `../my-texas-42-react-app/packages/functions/src/utils/lobby-utils.ts`
- Connect handler: `../my-texas-42-react-app/packages/functions/src/websockets/connect.ts`
- Play-turn handler: `../my-texas-42-react-app/packages/functions/src/websockets/play-turn.ts`

The legacy code has no tests, so rule edges (Sevens ties, doubles-own-suit cross-suit interactions, mark-bid sequences) need to be confirmed against game references / play sessions before you finalize a test fixture.
