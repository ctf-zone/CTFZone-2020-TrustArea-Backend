package models

import (
	"api/app/db"
	"api/app/utils"
	"api/config"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func CreateUser(username, firstName, lastName string) (string, error) {
	var uuid string

	err := db.DB.QueryRow("INSERT into users (username, first_name, last_name) VALUES ($1, $2, $3) RETURNING id",
		username, firstName, lastName).Scan(&uuid)

	pgerr, ok := err.(pgx.PgError)

	if err != nil {
		if ok == true && pgerr.Code == "23505" {
			err = utils.UserAlreadyExistsError
		}
		log.WithFields(log.Fields{
			"username": username,
			"firstName": firstName,
			"lastName": lastName,
		}).Error(err)
		return "", err
	}

	log.WithFields(log.Fields{
		"username": username,
		"firstName": firstName,
		"lastName": lastName,
		"uuid": uuid,
	}).Info("New user!")

	return uuid, nil
}

func CreateSession(refreshToken string) (string, error) {
	var err error

	var tokensCount int
	err = db.DB.QueryRow("SELECT count(*) FROM USERS WHERE id = $1", refreshToken).Scan(&tokensCount)

	pgerr, ok := err.(pgx.PgError)

	if err != nil {
		if ok == true && pgerr.Code == "22P02" && strings.HasPrefix(pgerr.Message, "invalid input syntax for type uuid") {
			err = utils.RefreshTokenNotExistsError
		}
		log.WithFields(log.Fields{
			"refreshToken": refreshToken,
		}).Error(err)
		return "", err
	}

	if tokensCount != 1 {
		return "", utils.RefreshTokenNotExistsError
	}

	token := make([]byte, 8)
	_,_ = rand.Read(token)
	session := sha256.Sum256([]byte(token))
	sessionToken := fmt.Sprintf("%x", session)
	err = db.Redis.Set(sessionToken, refreshToken, config.Config.Redis.SessionTTL * time.Second).Err()

	if err != nil {
		log.WithFields(log.Fields{
			"sessionToken": sessionToken,
			"refreshToken": refreshToken,
			"sessionTTL": config.Config.Redis.SessionTTL,
		}).Error(err)
		return "", err
	}

	log.WithFields(log.Fields{
		"sessionToken": sessionToken,
		"refreshToken": refreshToken,
		"sessionTTL": config.Config.Redis.SessionTTL,
	}).Info("New session!")

	return sessionToken, nil
}

func CheckSession(session string) (string, error) {
	val, err := db.Redis.Get(session).Result()

	if err != nil && err != redis.Nil {
		log.WithFields(log.Fields{
			"sessionToken": session,
		}).Error(err)
		return "", err
	}

	return val, nil
}
