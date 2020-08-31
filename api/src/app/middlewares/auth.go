package middlewares

import (
	"api/app/models"
	"api/app/utils"
	"context"
	"encoding/json"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "",
			"errors": []string{},
			"data": map[string]string{},
		}

		handlersWithoutAuth := []string {
			"/health",
			"/logging",
			"/auth/register",
			"/auth/session",
		}

		if utils.ArrayContainsString(handlersWithoutAuth, r.RequestURI) == true {
			next.ServeHTTP(w, r)
		} else {
			sessionToken := r.Header.Get("Authorization")

			if sessionToken == "" {
				w.WriteHeader(http.StatusForbidden)
				response["message"] = "Missing Session Token"
				json.NewEncoder(w).Encode(response)
				return
			}

			uuid, _ := models.CheckSession(sessionToken)

			if uuid == "" {
				w.WriteHeader(http.StatusForbidden)
				response["message"] = "Session Token is not found"
				json.NewEncoder(w).Encode(response)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user_uuid", uuid))
			next.ServeHTTP(w, r)
		}
	})
}