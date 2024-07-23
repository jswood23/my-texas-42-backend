package sql_scripts

import "fmt"

func RemoveFriendRequest(senderUsername string, receiverUsername string) string {
	return fmt.Sprintf(`
WITH UserIDs AS (
	SELECT UserID, 1 AS x FROM public.Users WHERE Username = '%s' UNION ALL
	SELECT UserID, 2 AS x FROM public.Users WHERE Username = '%s'
)
DELETE FROM public.FriendRequests
WHERE senderuserid = ( SELECT UserID FROM UserIDs WHERE x = 1 ) AND receiveruserid = ( SELECT UserID FROM UserIDs WHERE x = 2 )
OR senderuserid = ( SELECT UserID FROM UserIDs WHERE x = 2 ) AND receiveruserid = ( SELECT UserID FROM UserIDs WHERE x = 1 );
`, senderUsername, receiverUsername)
}
