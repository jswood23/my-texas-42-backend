package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
	"strconv"
)

func SaveChatMessage(matchId int, username string, message string) string {
	sanitizedMatchId := util.Sanitize(strconv.Itoa(matchId))
	sanitizedUsername := util.Sanitize(username)
	sanitizedMessage := util.Sanitize(message)
	return fmt.Sprintf(`
INSERT INTO public.ChatMessageArchive (MatchID, Username, Message)
VALUES ('%s', '%s', '%s');
`, sanitizedMatchId, sanitizedUsername, sanitizedMessage)
}
