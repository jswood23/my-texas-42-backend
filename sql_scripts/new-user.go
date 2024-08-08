package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func NewUser(email string, username string) string {
	sanitizedEmail, sanitizedUsername := util.Sanitize(email), util.Sanitize(username)

	return fmt.Sprintf(`
INSERT INTO public.Users (Email, Username, IsAdmin, DisplayName)
VALUES ('%s', '%s', FALSE);
INSERT INTO public.UserStats (UserId)
SELECT UserID FROM public.Users
WHERE Username = '%s';
`, sanitizedEmail, sanitizedUsername, sanitizedUsername)
}
