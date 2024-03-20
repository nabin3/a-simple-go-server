package main

import (
	"fmt"
	"net/http"
)

// Defining handler for checking readiness of server throgh "/healthz"
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("%s\n", http.StatusText(http.StatusOK))))
}
