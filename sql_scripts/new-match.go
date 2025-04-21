package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/util"
)

func NewMatch(matchName string, matchPrivacy string, rules string, team1Player1 string) string {
	sanitizedMatchName := util.Sanitize(matchName)
	sanitizedMatchPrivacy := util.Sanitize(matchPrivacy)
	sanitizedRules := util.Sanitize(rules)
	sanitizedTeam1Player1 := util.Sanitize(team1Player1)
	return fmt.Sprintf(`
INSERT INTO public.MatchArchive
(MatchName, MatchPrivacy, Rules, Team1Player1)
VALUES ('%s', '%s', '%s', '%s');
SELECT MatchID
FROM public.MatchArchive
ORDER BY MatchID DESC LIMIT 1;
`, sanitizedMatchName, sanitizedMatchPrivacy, sanitizedRules, sanitizedTeam1Player1)
}
