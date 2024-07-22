package sql_scripts

import "fmt"

func CreateUser(email string, username string) string {
	return fmt.Sprintf(`
INSERT INTO public.Users (Email, Username, IsAdmin, DisplayName)
VALUES ('%s', '%s', FALSE);
INSERT INTO public.UserStats (UserId)
SELECT UserID FROM public.Users
WHERE Username = '%s';
`, email, username, username)
}
