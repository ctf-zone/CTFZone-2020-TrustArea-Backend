package handlers

import (
	"api/app/utils"
	"encoding/json"
	"net/http"
)

func PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "",
		"errors": []string{},
		"data": map[string]string{},
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response["message"] = utils.PageNotFoundError.Error()
	json.NewEncoder(w).Encode(response)
}
