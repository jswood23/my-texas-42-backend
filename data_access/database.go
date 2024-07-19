package data_access

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
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

	dbUrl = os.Getenv("HOST_NAME")
	dbPort = os.Getenv("POSTGRES_PRODUCTION_PORT")
	dbName = os.Getenv("DB_NAME")
	dbUsername = os.Getenv("POSTGRES_USER")
	dbPassword = os.Getenv("POSTGRES_PASSWORD")

	if environment == "staging" {
		dbPort = os.Getenv("POSTGRES_STAGING_PORT")
	}

	if dbUrl == "" || dbPort == "" || dbName == "" || dbUsername == "" || dbPassword == "" {
		return errors.New(missingEnvVarsMessage)
	}

	return nil
}

func GetDBSession() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbUrl, dbPort, dbUsername, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}

func CheckDBConnection() error {
	db, err := GetDBSession()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
