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

type GetSessionRequest struct {
	RefreshToken string `json:"refresh_token" valid:"required,type(string),length(1|1024)"`
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var getSessionRequest GetSessionRequest
	var sessionToken string
	var err error
	var errors []string

	_ = json.NewDecoder(r.Body).Decode(&getSessionRequest)
	validationResult, validationErrors := govalidator.ValidateStruct(getSessionRequest)

	if validationResult != true {
		err = utils.ValidationError
		log.WithFields(log.Fields{
			"RefreshToken": getSessionRequest.RefreshToken,
		}).Error("RegisterUserHandler validation errors!")

		validationErrorsAsserted, _ := validationErrors.(govalidator.Errors)
		for _, e := range validationErrorsAsserted {
			eAsserted, _ := e.(govalidator.Error)
			variableName := eAsserted.Name
			errorMessage := eAsserted.Err.Error()
			errors = append(errors, fmt.Sprintf("%s - %s", variableName, strings.TrimSpace(errorMessage)))
		}

	} else {
		sessionToken, err = models.CreateSession(getSessionRequest.RefreshToken)
	}

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		response["message"] = ""
		response["data"] = map[string]string{
			"token_type": "session",
			"token": sessionToken,
		}
		json.NewEncoder(w).Encode(response)
	case utils.ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		response["message"] = utils.ValidationError.Error()
		response["errors"] = errors
		json.NewEncoder(w).Encode(response)
	case utils.RefreshTokenNotExistsError:
		w.WriteHeader(http.StatusNotFound)
		response["message"] = utils.RefreshTokenNotExistsError.Error()
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = utils.UndefinedError.Error()
		json.NewEncoder(w).Encode(response)
	}
}
