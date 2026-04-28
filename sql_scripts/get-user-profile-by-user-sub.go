package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func GetUserProfileByUserSub(username string, userSub string) string {
	sanitizedUsername, sanitizedUserSub := util.Sanitize(username), util.Sanitize(userSub)

	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE Username = '%s'
AND UserSub = '%s';
`, sanitizedUsername, sanitizedUserSub)
}

func GetUserByUserSub(userSub string) string {
	sanitizedUserSub := util.Sanitize(userSub)

	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE UserSub = '%s';
`, sanitizedUserSub)
}
