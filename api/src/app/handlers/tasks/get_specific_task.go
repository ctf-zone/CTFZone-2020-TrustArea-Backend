package tasks_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetSpecificTaskHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var err error
	var task *models.Task
	vars := mux.Vars(r)
	id := vars["id"]
	uuid := fmt.Sprintf("%v", r.Context().Value("user_uuid"))

	var idInt int
	if idInt, err = strconv.Atoi(id); err == nil {
		task, err = models.GetTask(idInt, uuid, false)
	} else {
		err = utils.UndefinedError
		log.WithFields(log.Fields{
			"TaskID": id,
		}).Error("Error while converting id to integer")
	}

	if task == nil {
		err = utils.TaskNotFoundError
	}

	switch err {
	case nil:
		response["data"] = task
		json.NewEncoder(w).Encode(response)
	case utils.TaskNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		response["message"] = utils.TaskNotFoundError.Error()
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = utils.UndefinedError.Error()
		json.NewEncoder(w).Encode(response)
	}
}
