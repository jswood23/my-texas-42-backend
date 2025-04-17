package logger

import (
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/system"
)

const (
	info models.LogLevel = iota
	warning
	errorLevel
	critical
)

func log(logLevel models.LogLevel, message string) {
	serviceName := "my-texas-42-backend-" + system.GetEnv()
	query := sql_scripts.InsertLog(serviceName, logLevel, message)
	_ = services.Execute(query)
}

func Info(message string) {
	log(info, message)
}

func Warning(message string) {
	log(warning, message)
}

func Error(message string) {
	log(errorLevel, message)
}

func Critical(message string) {
	defer println("CRITICAL: " + message)
	log(critical, message)
}
