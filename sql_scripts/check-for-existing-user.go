package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func CheckForExistingUser(username string, email string) string {
	sanitizedUsername, sanitizedEmail := util.Sanitize(username), util.Sanitize(email)

	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE Username = '%s'
OR Email = '%s';
`, sanitizedUsername, sanitizedEmail)
}
