package models

import (
	"api/app/db"
	log "github.com/sirupsen/logrus"
)

func SubmitSolution(taskID int, solution string, uuid string) (int, error) {
	var solutionID int

	err := db.DB.QueryRow("INSERT into solutions (task_id, solution, owner) VALUES ($1, $2, $3) RETURNING id",
		taskID, solution, uuid).Scan(&solutionID)

	if err != nil {
		log.WithFields(log.Fields{
			"taskID": taskID,
			"solution": solution,
			"uuid": uuid,
		}).Error(err)
		return 0, err
	}

	return solutionID, nil
}