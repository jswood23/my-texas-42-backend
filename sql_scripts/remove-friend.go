package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func RemoveFriend(senderUsername string, receiverUsername string) string {
	sanitizedSenderUsername, sanitizedReceiverUsername := util.Sanitize(senderUsername), util.Sanitize(receiverUsername)

	return fmt.Sprintf(`
WITH UserIDs AS (
	SELECT UserID, 1 AS x FROM public.Users WHERE Username = '%s' UNION ALL
	SELECT UserID, 2 AS x FROM public.Users WHERE Username = '%s'
)
DELETE FROM public.Friends
WHERE User1ID = ( SELECT UserID FROM UserIDs WHERE x = 1 ) AND User2ID = ( SELECT UserID FROM UserIDs WHERE x = 2 )
OR User2ID = ( SELECT UserID FROM UserIDs WHERE x = 2 ) AND User1ID = ( SELECT UserID FROM UserIDs WHERE x = 1 );
`, sanitizedSenderUsername, sanitizedReceiverUsername)
}
