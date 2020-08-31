package tasks_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type CreateTaskRequest struct {
	Description string `json:"description" valid:"required,type(string),length(1|1024)"`
	Challenge string `json:"challenge" valid:"required,type(string),length(1|1024)"`
	Reward string `json:"reward" valid:"required,type(string),length(1|1024)"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var createTaskRequest CreateTaskRequest
	var err error
	var taskID int
	var errors []string

	_ = json.NewDecoder(r.Body).Decode(&createTaskRequest)
	validationResult, validationErrors := govalidator.ValidateStruct(createTaskRequest)

	if validationResult != true {
		err = utils.ValidationError
		log.WithFields(log.Fields{
			"Description": createTaskRequest.Description,
			"Challenge": createTaskRequest.Challenge,
			"Reward": createTaskRequest.Reward,
		}).Error("CreateTaskHandler validation errors!")

		validationErrorsAsserted, _ := validationErrors.(govalidator.Errors)
		for _, e := range validationErrorsAsserted {
			eAsserted, _ := e.(govalidator.Error)
			variableName := eAsserted.Name
			errorMessage := eAsserted.Err.Error()
			errors = append(errors, fmt.Sprintf("%s - %s", variableName, strings.TrimSpace(errorMessage)))
		}

	} else {
		uuid := fmt.Sprintf("%v", r.Context().Value("user_uuid"))
		taskID, err = models.CreateTask(createTaskRequest.Description, createTaskRequest.Challenge,
			createTaskRequest.Reward, uuid)
	}

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		response["message"] = ""
		response["data"] = map[string]int{
			"task_id": taskID,
		}
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