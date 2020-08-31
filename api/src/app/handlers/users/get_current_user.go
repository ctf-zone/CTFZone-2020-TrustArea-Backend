package users_handler

import (
	"api/app/models"
	"api/app/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}

	var currentUser *models.CurrentUser
	var err error

	uuid := fmt.Sprintf("%v", r.Context().Value("user_uuid"))
	currentUser, err = models.GetCurrentUser(uuid)

	switch err {
	case nil:
		response["data"] = currentUser
		json.NewEncoder(w).Encode(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = utils.UndefinedError.Error()
		json.NewEncoder(w).Encode(response)
	}
}