package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func NewUser(email string, username string, usersub string) string {
	sanitizedEmail, sanitizedUsername, sanitizedUserSub := util.Sanitize(email), util.Sanitize(username), util.Sanitize(usersub)

	return fmt.Sprintf(`
INSERT INTO public.Users (Email, Username, IsAdmin, DisplayName, UserSub)
VALUES ('%s', '%s', FALSE, '%s');
INSERT INTO public.UserStats (UserId)
SELECT UserID FROM public.Users
WHERE Username = '%s';
`, sanitizedEmail, sanitizedUsername, sanitizedUsername, sanitizedUserSub)
}
