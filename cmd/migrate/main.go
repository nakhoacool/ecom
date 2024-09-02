package main

import (
	"ecom/config"
	"ecom/db"
	"errors"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	db, err := db.NewMySQLStorage(mysqlCfg.Config{
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

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "mysql", driver)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}
}
