package models

import (
	"api/app/db"
	log "github.com/sirupsen/logrus"
)

func SaveLog(logType, logMessage string) (int, error) {
	var id int

	err := db.DB.QueryRow("INSERT into logs (log_type, log_message) VALUES ($1, $2) RETURNING id",
		logType, logMessage).Scan(&id)

	if err != nil {
		log.WithFields(log.Fields{
			"logType": logType,
			"logMessage": logMessage,
		}).Error(err)
	}

	return id, nil
}
