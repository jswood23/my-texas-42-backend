package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func ChangeDisplayName(newDisplayName string, username string) string {
	sanitizedDisplayName, sanitizedUsername := util.Sanitize(newDisplayName), util.Sanitize(username)

	return fmt.Sprintf(`
UPDATE users
SET displayname = '%s'
WHERE username = '%s';
`, sanitizedDisplayName, sanitizedUsername)
}
