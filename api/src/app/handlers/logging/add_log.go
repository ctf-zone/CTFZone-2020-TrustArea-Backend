package logging_handler

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

type AddLogRequest struct {
	LogType string `json:"log_type" valid:"required,type(string),length(1|32)"`
	LogMessage string `json:"log_message" valid:"required,type(string),length(1|1024)"`
}

func AddLogHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var addLogRequest AddLogRequest
	var id int
	var err error
	var errors []string

	_ = json.NewDecoder(r.Body).Decode(&addLogRequest)
	validationResult, validationErrors := govalidator.ValidateStruct(addLogRequest)

	if validationResult != true {
		err = utils.ValidationError
		log.WithFields(log.Fields{
			"logType": addLogRequest.LogType,
			"logMessage": addLogRequest.LogMessage,
		}).Error("AddLogHandler validation errors!")

		validationErrorsAsserted, _ := validationErrors.(govalidator.Errors)
		for _, e := range validationErrorsAsserted {
			eAsserted, _ := e.(govalidator.Error)
			variableName := eAsserted.Name
			errorMessage := eAsserted.Err.Error()
			errors = append(errors, fmt.Sprintf("%s - %s", variableName, strings.TrimSpace(errorMessage)))
		}

	} else {
		id, err = models.SaveLog(addLogRequest.LogType, addLogRequest.LogMessage)
	}

	switch err {
	case nil:
		log.WithFields(log.Fields{
			"LogID": id,
			"LogType": addLogRequest.LogType,
			"LogMessage": addLogRequest.LogMessage,
		}).Info("New log record!")
		w.WriteHeader(http.StatusOK)
		response["message"] = "Log saved!"
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