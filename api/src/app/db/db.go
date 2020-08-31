package db

import (
	"api/config"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	log "github.com/sirupsen/logrus"
	"os"
)

var DB *sql.DB

func DBInit() {
	var err error

	addr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Config.DB.User, config.Config.DB.Pass,
		config.Config.DB.Host, config.Config.DB.Port, config.Config.DB.DBName)
	DB, err = sql.Open("pgx", addr)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = DB.Ping()

	if err != nil {
		log.Error("Error while connecting to PostgreSQL")
		log.Error(err)
		os.Exit(1)
	}

	log.Infof("Connected to PostgreSQL: %s", addr)
}