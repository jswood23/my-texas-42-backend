package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"my-texas-42-backend/system"
	"reflect"
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

// Query Note: T must be a struct with the same properties as the query returns in the same order
func Query[T any](query string) ([]T, error) {
	db, err := GetDBSession()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var rows []T
	for result.Next() {
		var row T
		val := reflect.ValueOf(&row).Elem()
		var scanArgs []interface{}
		for i := 0; i < val.NumField(); i++ {
			scanArgs = append(scanArgs, val.Field(i).Addr().Interface())
		}

		err = result.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}
