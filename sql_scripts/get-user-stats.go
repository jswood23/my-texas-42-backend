package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func GetUserStats(username string) string {
	sanitizedUsername := util.Sanitize(username)
	return fmt.Sprintf(`
WITH UserID AS ( SELECT userid as id FROM users WHERE username = '%s' )
SELECT * FROM userstats WHERE userid = (SELECT id FROM UserID);
`, sanitizedUsername)
}
