package common

import (
	"fmt"
	"matsya/src/config"

	"github.com/jmoiron/sqlx"
)

func GetConnection() (*sqlx.DB, error) {
	dbConfig := config.GetConfigMap("database")

	dialect := dbConfig["dialect"]
	user := dbConfig["user"]
	dbName := dbConfig["database_name"]
	sslMode := dbConfig["sslmode"]

	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=%s", user, dbName, sslMode)

	db, err := sqlx.Connect(dialect, connectionString)
	if err != nil {
		fmt.Printf(">> error: %s", err.Error())
		return nil, err
	}

	return db, nil
}
