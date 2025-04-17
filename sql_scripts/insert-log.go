package sql_scripts

import (
	"fmt"
	"my-texas-42-backend/models"
	"my-texas-42-backend/util"
	"strconv"
)

func InsertLog(serviceName string, logLevel models.LogLevel, message string) string {
	sanitizedServiceName, sanitizedLogLevel, sanitizedMessage := util.Sanitize(serviceName), strconv.Itoa(int(logLevel)), util.Sanitize(message)

	return fmt.Sprintf(`
INSERT INTO logs (servicename, loglevel, message)
VALUES ('%s', '%s', '%s');
`, sanitizedServiceName, sanitizedLogLevel, sanitizedMessage)
}
