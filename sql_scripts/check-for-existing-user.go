package sql_scripts

import "fmt"

func CheckForExistingUser(username string, email string) string {
	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE Username = '%s'
OR Email = '%s';
`, username, email)
}
