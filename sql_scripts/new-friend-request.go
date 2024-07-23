package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func NewFriendRequest(senderUsername string, receiverUsername string) string {
	sanitizedSenderUsername, sanitizedReceiverUsername := util.Sanitize(senderUsername), util.Sanitize(receiverUsername)

	return fmt.Sprintf(`
INSERT INTO public.FriendRequests (senderuserid, receiveruserid)
VALUES (
        	(SELECT userid FROM users WHERE username = '%s'),
        	(SELECT userid FROM users WHERE username = '%s')
);
`, sanitizedSenderUsername, sanitizedReceiverUsername)
}
