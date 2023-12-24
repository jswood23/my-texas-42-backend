package data_access

import (
	"errors"
	"os"
)

var dbAddress = ""
var dbUsername = ""
var dbPassword = ""

func Initialize() error {
	if len(os.Args) < 4 {
		return errors.New("not enough program arguments")
	}

	dbAddress = os.Args[1]
	dbUsername = os.Args[2]
	dbPassword = os.Args[3]

	return nil
}

func GetDBLogin() (string, string, string) {
	return dbAddress, dbUsername, dbPassword
}
