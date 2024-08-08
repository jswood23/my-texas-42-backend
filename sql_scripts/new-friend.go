package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func NewFriend(user1Username string, user2Username string) string {
	sanitizedUser1Username, sanitizedUser2Username := util.Sanitize(user1Username), util.Sanitize(user2Username)

	return fmt.Sprintf(`
WITH UserIDs AS (
	SELECT UserID, 1 AS x FROM public.Users WHERE Username = '%s' UNION ALL
	SELECT UserID, 2 AS x FROM public.Users WHERE Username = '%s'
)
INSERT INTO friends (user1id, user2id)
VALUES ((SELECT UserID FROM UserIDs WHERE x = 1), (SELECT UserID FROM UserIDs WHERE x = 2));
`, sanitizedUser1Username, sanitizedUser2Username)
}
