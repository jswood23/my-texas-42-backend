package sql_scripts

import "fmt"

func GetUserProfileByUsername(username string) string {
	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE Username = '%s';
`, username)
}
