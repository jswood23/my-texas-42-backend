package data_access

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"my-texas-42-backend/system"
)

func GetDBSession() (*sql.DB, error) {
	dbUrl := system.GetDBUrl()
	dbPort := system.GetDBPort()
	dbUsername := system.GetDBUsername()
	dbPassword := system.GetDBPassword()
	dbName := system.GetDBName()

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
