package auth_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"strings"

	// "log"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RegisterUserRequest struct {
	Username string `json:"username" valid:"required,type(string),length(1|1024)"`
	FirstName string `json:"first_name" valid:"required,type(string),length(1|1024)"`
	LastName string `json:"last_name" valid:"required,type(string),length(1|1024)"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var registerUserRequest RegisterUserRequest
	var uuid string
	var err error
	var errors []string

	_ = json.NewDecoder(r.Body).Decode(&registerUserRequest)
	validationResult, validationErrors := govalidator.ValidateStruct(registerUserRequest)

	if validationResult != true {
		err = utils.ValidationError
		log.WithFields(log.Fields{
			"Username": registerUserRequest.Username,
			"FirstName": registerUserRequest.FirstName,
			"LastName": registerUserRequest.LastName,
		}).Error("RegisterUserHandler validation errors!")

		validationErrorsAsserted, _ := validationErrors.(govalidator.Errors)
		for _, e := range validationErrorsAsserted {
			eAsserted, _ := e.(govalidator.Error)
			variableName := eAsserted.Name
			errorMessage := eAsserted.Err.Error()
			errors = append(errors, fmt.Sprintf("%s - %s", variableName, strings.TrimSpace(errorMessage)))
		}

	} else {
		uuid, err = models.CreateUser(registerUserRequest.Username, registerUserRequest.FirstName,
			registerUserRequest.LastName)
	}

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		response["message"] = ""
		response["data"] = map[string]string{
			"token_type": "refresh",
			"token": uuid,
		}
		json.NewEncoder(w).Encode(response)
	case utils.UserAlreadyExistsError:
		w.WriteHeader(http.StatusConflict)
		response["message"] = utils.UserAlreadyExistsError.Error()
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