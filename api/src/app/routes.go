package app

import (
	"api/app/handlers"
	"api/app/handlers/auth"
	"api/app/handlers/logging"
	solutions_handler "api/app/handlers/solutions"
	tasks_handler "api/app/handlers/tasks"
	users_handler "api/app/handlers/users"
	"api/app/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func GenerateRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(middlewares.СommonMiddleware)
	r.Use(middlewares.AuthMiddleware)

	// Health Check
	r.HandleFunc("/health", handlers.HealthCheckPage).Methods(http.MethodGet)

	// Auth
	r.HandleFunc("/auth/register", auth_handler.RegisterUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/auth/session", auth_handler.GetSessionHandler).Methods(http.MethodPost)

	// Logging
	r.HandleFunc("/logging", logging_handler.AddLogHandler).Methods(http.MethodPost)

	// Tasks
	r.HandleFunc("/tasks", tasks_handler.CreateTaskHandler).Methods(http.MethodPost)
	r.HandleFunc("/tasks", tasks_handler.GetTasksHandler).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id:[0-9]+}", tasks_handler.UpdateTaskHandler).Methods(http.MethodPut)
	r.HandleFunc("/tasks/{id:[0-9]+}", tasks_handler.GetSpecificTaskHandler).Methods(http.MethodGet)

	// Users
	r.HandleFunc("/users", users_handler.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/users/me", users_handler.GetCurrentUser).Methods(http.MethodGet)

	// Solutions
	r.HandleFunc("/solutions", solutions_handler.SubmitSolutionHandler).Methods(http.MethodPost)

	r.NotFoundHandler = http.HandlerFunc(handlers.PageNotFoundHandler)

	return r
}
