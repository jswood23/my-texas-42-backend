package system

import (
	"errors"
	"os"
)

var environment = ""
var userPoolName = ""
var userPoolAppKey = ""
var dbUrl = ""
var dbPort = ""
var dbName = ""
var dbUsername = ""
var dbPassword = ""
var missingDBVarsMessage = "one or more database environment variables are not set"
var missingUserPoolVarsMessage = "one or more user pool environment variables are not set"

func Initialize() error {
	environment = os.Getenv("ENVIRONMENT")
	if environment == "" {
		return errors.New(missingDBVarsMessage)
	}

	err := initializeDB()
	if err != nil {
		return err
	}

	err = initializeUserPool()
	if err != nil {
		return err
	}

	return nil
}

func initializeDB() error {
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
		return errors.New(missingDBVarsMessage)
	}

	return nil
}

func initializeUserPool() error {
	userPoolName = os.Getenv("STAGING_USER_POOL_NAME")
	userPoolAppKey = os.Getenv("STAGING_USER_POOL_APP_KEY")

	if environment == "production" {
		userPoolName = os.Getenv("PRODUCTION_USER_POOL_NAME")
		userPoolAppKey = os.Getenv("PRODUCTION_USER_POOL_APP_KEY")
	}

	if userPoolName == "" || userPoolAppKey == "" {
		return errors.New(missingUserPoolVarsMessage)
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

func GetUserPoolName() string {
	return userPoolName
}

func GetUserPoolAppKey() string {
	return userPoolAppKey
}
