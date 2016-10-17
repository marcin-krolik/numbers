package mgmt

import (
	"fmt"
	"net/http"
)

// rootHandler for "/" path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server version %d \n", ServerVersion)
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

