package data_access

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

var dbUrl = ""
var dbPort = ""
var dbName = ""
var dbUsername = ""
var dbPassword = ""
var authSecretKey = ""

func Initialize() error {
	if len(os.Args) < 7 {
		println("args:")
		for _, arg := range os.Args {
			println(arg)
		}

		return errors.New("not enough program arguments")
	}

	dbUrl = os.Args[1]
	dbPort = os.Args[2]
	dbName = os.Args[3]
	dbUsername = os.Args[4]
	dbPassword = os.Args[5]
	authSecretKey = os.Args[6]

	return nil
}

func GetDBSession() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbUrl, dbPort, dbUsername, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}

func GetAuthSecretKey() string {
	return authSecretKey
}
