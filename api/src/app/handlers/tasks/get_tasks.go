package tasks_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var tasks []*models.Task
	var err error
	var limit, offset int
	var username string

	limitArr := r.URL.Query()["limit"]
	if len(limitArr) > 0 {
		var errLocal error
		limit, errLocal = strconv.Atoi(limitArr[0])
		if errLocal != nil {
			limit = 100
		}
	} else {
		limit = 100
	}

	offsetArr := r.URL.Query()["offset"]
	if len(offsetArr) > 0 {
		var errLocal error
		offset, errLocal = strconv.Atoi(offsetArr[0])
		if errLocal != nil {
			offset = 0
		}
	} else {
		offset = 0
	}

	usernameArr := r.URL.Query()["username"]
	if len(usernameArr) > 0 {
		username = usernameArr[0]
	} else {
		username = ""
	}

	uuid := fmt.Sprintf("%v", r.Context().Value("user_uuid"))
	tasks, err = models.GetTasks(uuid, offset, limit, username)

	switch err {
	case nil:
		response["data"] = tasks
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = utils.UndefinedError.Error()
		json.NewEncoder(w).Encode(response)
	}
}
