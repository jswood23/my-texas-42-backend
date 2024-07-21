package system

import (
	"errors"
	"os"
)

var dbUrl = ""
var dbPort = ""
var dbName = ""
var dbUsername = ""
var dbPassword = ""
var environment = ""
var missingEnvVarsMessage = "one or more environment variables are not set"

func Initialize() error {
	environment = os.Getenv("ENVIRONMENT")
	if environment == "" {
		return errors.New(missingEnvVarsMessage)
	}

	// this is the same for both staging and production because we are going to the same port within a different container
	dbUrl = os.Getenv("POSTGRES_PRODUCTION_HOST_NAME")

	dbPort = os.Getenv("POSTGRES_PRODUCTION_PORT")
	dbName = os.Getenv("DB_NAME")
	dbUsername = os.Getenv("POSTGRES_USER")
	dbPassword = os.Getenv("POSTGRES_PASSWORD")

	if environment == "staging" {
		dbUrl = os.Getenv("POSTGRES_STAGING_HOST_NAME")
	}

	if dbUrl == "" || dbPort == "" || dbName == "" || dbUsername == "" || dbPassword == "" {
		return errors.New(missingEnvVarsMessage)
	}

	return nil
}

func GetEnv() string {
	return environment
}

func GetDBUrl() string {
	return dbUrl
}

func GetDBPort() string {
	return dbPort
}

func GetDBName() string {
	return dbName
}

func GetDBUsername() string {
	return dbUsername
}

func GetDBPassword() string {
	return dbPassword
}
