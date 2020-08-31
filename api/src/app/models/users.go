package models

import (
	"api/app/db"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type User struct {
	ID string
	Username string
	FirstName string
	LastName string
	Date time.Time
}

type CurrentUser struct {
	ID string
	Username string
	FirstName string
	LastName string
	Date time.Time
	// Tasks[]
}

func GetUsers(limit string) ([]*User, error) {
	users := make([]*User, 0)

	var SQLRequest = `
		SELECT '', username, first_name, last_name, date
		FROM users
		ORDER BY date DESC
		LIMIT $1
	`

	rows, err := db.DB.Query(strings.Replace(SQLRequest, "$1", limit, 1))

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Date)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetSpecificUser(username string) (*User, error) {

	user := new(User)

	var SQLRequest = `
		SELECT id, username, first_name, last_name, date
		FROM users
		WHERE username = $1
	`

	err := db.DB.QueryRow(SQLRequest, username).Scan(&user.ID, &user.Username, &user.FirstName,
		&user.LastName, &user.Date)

	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
		}).Error(err)
		return nil, err
	}

	return user, nil
}

func GetCurrentUser(uuid string) (*CurrentUser, error) {

	currentUser := new(CurrentUser)

	var SQLRequest = `
		SELECT id, username, first_name, last_name, date
		FROM users
		WHERE id = $1
	`

	err := db.DB.QueryRow(SQLRequest, uuid).Scan(&currentUser.ID, &currentUser.Username, &currentUser.FirstName,
		&currentUser.LastName, &currentUser.Date)

	if err != nil {
		log.WithFields(log.Fields{
			"uuid": uuid,
		}).Error(err)
		return nil, err
	}

	return currentUser, nil

}