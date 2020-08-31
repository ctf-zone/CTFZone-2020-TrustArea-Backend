package tasks_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type UpdateTaskRequest struct {
	Description string `json:"description" valid:"type(string)"`
	Challenge string `json:"challenge" valid:"type(string)"`
	Reward string `json:"reward" valid:"type(string)"`
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var err error
	vars := mux.Vars(r)
	id := vars["id"]
	var updateTaskRequest UpdateTaskRequest
	var errors []string

	_ = json.NewDecoder(r.Body).Decode(&updateTaskRequest)
	validationResult, validationErrors := govalidator.ValidateStruct(updateTaskRequest)

	if validationResult != true {
		err = utils.ValidationError
		log.WithFields(log.Fields{
			"Description": updateTaskRequest.Description,
			"Challenge": updateTaskRequest.Challenge,
			"Reward": updateTaskRequest.Reward,
		}).Error("UpdateTaskRequest validation errors!")

		validationErrorsAsserted, _ := validationErrors.(govalidator.Errors)
		for _, e := range validationErrorsAsserted {
			eAsserted, _ := e.(govalidator.Error)
			variableName := eAsserted.Name
			errorMessage := eAsserted.Err.Error()
			errors = append(errors, fmt.Sprintf("%s - %s", variableName, strings.TrimSpace(errorMessage)))
		}
	} else {
		var idInt int
		if idInt, err = strconv.Atoi(id); err == nil {
			err = models.UpdateTask(idInt, updateTaskRequest.Description, updateTaskRequest.Challenge, updateTaskRequest.Reward)
		} else {
			err = utils.UndefinedError
			log.WithFields(log.Fields{
				"TaskID": id,
			}).Error("Error while converting id to integer")
		}
	}

	switch err {
	case nil:
		response["message"] = "OK"
		json.NewEncoder(w).Encode(response)
	case utils.ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		response["message"] = utils.ValidationError.Error()
		response["errors"] = errors
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = utils.UndefinedError.Error()
		json.NewEncoder(w).Encode(response)
	}
}