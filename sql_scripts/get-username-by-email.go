package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func GetUsernameByEmail(email string) string {
	sanitizedEmail := util.Sanitize(email)

	return fmt.Sprintf(`
SELECT username
FROM users
WHERE email = '%s';
`, sanitizedEmail)
}
