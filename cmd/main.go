package main

import (
	"database/sql"
	"ecom/cmd/api"
	"ecom/config"
	"ecom/db"
	"github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.ENV.DBUser,
		Passwd:               config.ENV.DBPassword,
		Addr:                 config.ENV.DBAddress,
		DBName:               config.ENV.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected")
}
