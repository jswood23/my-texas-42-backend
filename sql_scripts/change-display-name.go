package sql_scripts

import "fmt"

func ChangeDisplayName(newDisplayName string, username string) string {
	return fmt.Sprintf(`
UPDATE users
SET displayname = '%s'
WHERE username = '%s';
`, newDisplayName, username)
}
