package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func CheckForExistingFriend(user1Username string, user2Username string) string {
	sanitizedUser1Username, sanitizedUser2Username := util.Sanitize(user1Username), util.Sanitize(user2Username)

	return fmt.Sprintf(`
WITH UserIDs AS (
	SELECT UserID, 1 AS x FROM public.Users WHERE Username = '%s' UNION ALL
	SELECT UserID, 2 AS x FROM public.Users WHERE Username = '%s'
)
SELECT * FROM public.Friends
WHERE user1id = ( SELECT UserID FROM UserIDs WHERE x = 1 ) AND user2id = ( SELECT UserID FROM UserIDs WHERE x = 2 )
OR user1id = ( SELECT UserID FROM UserIDs WHERE x = 2 ) AND user2id = ( SELECT UserID FROM UserIDs WHERE x = 1 );
`, sanitizedUser1Username, sanitizedUser2Username)
}
