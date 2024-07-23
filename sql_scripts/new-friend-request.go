package sql_scripts

import "fmt"

func NewFriendRequest(senderUsername string, receiverUsername string) string {
	return fmt.Sprintf(`
INSERT INTO public.FriendRequests (senderuserid, receiveruserid)
VALUES (
        	(SELECT userid FROM users WHERE username = '%s'),
        	(SELECT userid FROM users WHERE username = '%s')
);
`, senderUsername, receiverUsername)
}
