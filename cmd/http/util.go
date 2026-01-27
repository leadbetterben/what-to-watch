package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	fmt.Printf("%s\n", err)
}

func writeMethodError(w http.ResponseWriter, allowedMethod string) {
	writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("only %s is allowed", allowedMethod))
}
