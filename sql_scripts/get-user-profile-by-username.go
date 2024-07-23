package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func GetUserProfileByUsername(username string) string {
	sanitizedUsername := util.Sanitize(username)

	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE Username = '%s';
`, sanitizedUsername)
}
