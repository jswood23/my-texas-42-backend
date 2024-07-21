package sql_scripts

import "fmt"

func GetUser(userid int) string {
	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE UserID = %d;
`, userid)
}
