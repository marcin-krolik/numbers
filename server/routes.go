package server

import (
	"fmt"
	"net/http"
	"encoding/json"
)

// rootHandler for "/" path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	respJson, _ := json.Marshal(map[string][]string{"endpoints": []string{"/", "numbers"}})
	// write headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// write response
	w.Write(respJson)
}

// numbersHandler for "/numbers" path
func numbersHandler(w http.ResponseWriter, r *http.Request) {
	respJson := processQuery(r.URL.Query())
	LogDebug(fmt.Sprintf("json response %s\n", string(respJson)))

	// write headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// write response
	w.Write(respJson)
}
