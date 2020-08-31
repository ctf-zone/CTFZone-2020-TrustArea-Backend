package models

import (
	"api/app/db"
	"api/app/utils"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

type Task struct {
	ID int
	Description string
	Challenge string
	Reward string
	Owner string
	Solved bool
	Date time.Time
}

func CreateTask(description, challenge, reward, uuid string) (int, error) {
	var taskID int

	var SQLRequest = "INSERT into tasks (description, challenge, reward, owner) VALUES ($1, $2, $3, $4) RETURNING id"

	err := db.DB.QueryRow(SQLRequest, description, challenge, reward, uuid).Scan(&taskID)

	if err != nil {
		log.WithFields(log.Fields{
			"description": description,
			"challenge": challenge,
			"reward": reward,
			"uuid": uuid,
		}).Error(err)
	}

	log.WithFields(log.Fields{
		"description": description,
		"challenge": challenge,
		"reward": reward,
		"uuid": uuid,
	}).Info("New task")

	return taskID, nil
}

func UpdateTask(id int, description, challenge, reward string) error {

	var SQLRequest = `
	UPDATE tasks SET
		description = CASE WHEN $1 <> '' THEN $1 ELSE description END,
		challenge = CASE WHEN $2 <> '' THEN $2 ELSE challenge END,
		reward = CASE WHEN $3 <> '' THEN $3 ELSE reward END
	WHERE id = $4
	`

	c, err := db.DB.Query(SQLRequest, description, challenge, reward, id)

	if err != nil {
		log.WithFields(log.Fields{
			"id": id,
			"description": description,
			"challenge": challenge,
			"reward": reward,
		}).Error(err)
		return err
	}

	defer c.Close()

	log.WithFields(log.Fields{
		"id": id,
		"description": description,
		"challenge": challenge,
		"reward": reward,
	}).Info("Task updated")

	return nil
}

func GetTask(id int, uuid string, forCheckSolution bool) (*Task, error) {
	task := new(Task)
	var SQLRequest string
	var err error

	if forCheckSolution == true {
		SQLRequest = `
			SELECT tasks.id, users.username, tasks.description, tasks.challenge, tasks.reward, tasks.solved, tasks.date 
			FROM tasks
			INNER JOIN users ON (tasks.owner=users.id)
			WHERE tasks.id=$1
		`
		err = db.DB.QueryRow(SQLRequest, id).Scan(&task.ID, &task.Owner, &task.Description, &task.Challenge,
			&task.Reward, &task.Solved, &task.Date)
	} else {
		SQLRequest = `
			SELECT tasks.id, users.username, tasks.description, tasks.challenge, 
			CASE WHEN users.id=$1 THEN tasks.reward ELSE '***' END as reward,  
			tasks.solved, tasks.date 
			FROM tasks
			INNER JOIN users ON (tasks.owner=users.id)
			WHERE tasks.id=$2
		`
		err = db.DB.QueryRow(SQLRequest, uuid, id).Scan(&task.ID, &task.Owner, &task.Description, &task.Challenge,
			&task.Reward, &task.Solved, &task.Date)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.TaskNotFoundError
		}
		log.WithFields(log.Fields{
			"TaskID": id,
		}).Error(err)
		return nil, err
	}

	return task, nil
}

func GetTasks(uuid string, offset, limit int, username string) ([]*Task, error) {
	tasks := make([]*Task, 0)

	var err error
	var rows *sql.Rows

	var SQLRequest = `
			SELECT tasks.id, users.username, tasks.description, tasks.challenge, 
			CASE WHEN users.id=$1 THEN tasks.reward ELSE '***' END as reward, 
			tasks.solved, tasks.date 
			FROM tasks 
			INNER JOIN users ON (tasks.owner=users.id)
	`
	if username != "" {
		SQLRequest += "WHERE users.username = $escape$" + username + "$escape$"
	}
	SQLRequest += `
			ORDER BY id DESC
			OFFSET $2 LIMIT $3
	`

	rows, err = db.DB.Query(SQLRequest, uuid, offset, limit)

	if err != nil {
		log.WithFields(log.Fields{
			"uuid": uuid,
		}).Error(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := new(Task)
		err := rows.Scan(&task.ID, &task.Owner, &task.Description, &task.Challenge, &task.Reward, &task.Solved, &task.Date)
		if err != nil {
			log.WithFields(log.Fields{
				"uuid": uuid,
			}).Error(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}