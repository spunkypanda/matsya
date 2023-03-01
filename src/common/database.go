package common

import (
	"fmt"
	"log"
	"matsya/src/config"

	"github.com/jmoiron/sqlx"
)

func GetConnection() *sqlx.DB {
	dbConfig := config.GetConfigMap("database")

	dialect := dbConfig["dialect"]
	user := dbConfig["user"]
	dbName := dbConfig["database_name"]
	sslMode := dbConfig["sslmode"]

	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=%s", user, dbName, sslMode)

	db, err := sqlx.Connect(dialect, connectionString)
	if err != nil {
		log.Fatal(err)
		fmt.Printf(">> error: %s", err.Error())
	}

	return db
}
