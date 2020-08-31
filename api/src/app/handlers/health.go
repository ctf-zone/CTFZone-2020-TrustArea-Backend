package handlers

import (
	"fmt"
	"net/http"
)

func HealthCheckPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
