package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/models"
)

func GetUser(userid models.UserID) string {
	return fmt.Sprintf(`
SELECT * FROM public.Users
WHERE UserID = %d;
`, userid)
}
