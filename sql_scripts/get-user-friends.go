package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func GetUserFriends(username string) string {
	sanitizedUsername := util.Sanitize(username)

	return fmt.Sprintf(`
WITH UserID AS ( SELECT userid as id FROM users WHERE username = '%s' )
SELECT
	username
FROM users
JOIN (
	SELECT user1id AS id FROM friends WHERE user2id = (SELECT id FROM UserID) UNION ALL
	SELECT user2id AS id FROM friends WHERE user1id = (SELECT id FROM UserID)
) ON users.userid = id;
`, sanitizedUsername)
}
