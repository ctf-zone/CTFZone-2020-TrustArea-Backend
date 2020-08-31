package solutions_handler

import (
	"api/app/models"
	"api/app/utils"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type SubmitSolutionRequest struct {
	TaskID int `json:"task_id" valid:"required,type(int)"`
	Solution string `json:"solution" valid:"required,type(string),length(1|1024)"`
}

func SubmitSolutionHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var submitSolutionRequest SubmitSolutionRequest
	var err error
	var task *models.Task
	var errors []string

	_ = json.NewDecoder(r.Body).Decode(&submitSolutionRequest)
	validationResult, validationErrors := govalidator.ValidateStruct(submitSolutionRequest)

	if validationResult != true {
		err = utils.ValidationError
		log.WithFields(log.Fields{
			"TaskID": submitSolutionRequest.TaskID,
			"Solution": submitSolutionRequest.Solution,
		}).Error("SubmitSolutionHandler validation errors!")

		validationErrorsAsserted, _ := validationErrors.(govalidator.Errors)
		for _, e := range validationErrorsAsserted {
			eAsserted, _ := e.(govalidator.Error)
			variableName := eAsserted.Name
			errorMessage := eAsserted.Err.Error()
			errors = append(errors, fmt.Sprintf("%s - %s", variableName, strings.TrimSpace(errorMessage)))
		}

	} else {
		uuid := fmt.Sprintf("%v", r.Context().Value("user_uuid"))
		hash := fmt.Sprintf("%x", sha1.Sum([]byte(submitSolutionRequest.Solution)))
		task, err = models.GetTask(submitSolutionRequest.TaskID, uuid, true)
		if err == nil {
			if hash == task.Challenge {
				_, err = models.SubmitSolution(task.ID, submitSolutionRequest.Solution, uuid)
			} else {
				err = utils.SolutionIsNotCorrectError
			}
		}
	}

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		response["message"] = ""
		response["data"] = map[string]string{
			"reward": task.Reward,
		}
		json.NewEncoder(w).Encode(response)
	case utils.SolutionIsNotCorrectError:
		w.WriteHeader(http.StatusNotAcceptable)
		response["message"] = utils.SolutionIsNotCorrectError.Error()
		json.NewEncoder(w).Encode(response)
	case utils.TaskNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		response["message"] = utils.TaskNotFoundError.Error()
		json.NewEncoder(w).Encode(response)
	case utils.ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		response["message"] = utils.ValidationError.Error()
		response["errors"] = errors
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = "Undefined error!"
		json.NewEncoder(w).Encode(response)
	}
}
