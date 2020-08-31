package users_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var users []*models.User
	var user *models.User
	var err error
	var limit string
	var username string

	usernameArr := r.URL.Query()["username"]
	if len(usernameArr) > 0 {
		username = usernameArr[0]
	} else {
		username = ""
	}

	limitArr := r.URL.Query()["limit"]
	if len(limitArr) > 0 {
		limit = limitArr[0]
	} else {
		limit = "100"
	}

	if username == "" {
		users, err = models.GetUsers(limit)
	} else {
		user, err = models.GetSpecificUser(username)
		if user == nil {
			err = utils.UserNotFoundError
		}
	}

	switch err {
	case nil:
		if user != nil {
			response["data"] = user
		} else {
			response["data"] = users
		}
		json.NewEncoder(w).Encode(response)
	case utils.UserNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		response["message"] = utils.UserNotFoundError.Error()
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = utils.UndefinedError.Error()
		json.NewEncoder(w).Encode(response)
	}
}
