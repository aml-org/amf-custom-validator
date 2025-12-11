package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aml-org/amf-custom-validator/internal/validator"
	"github.com/gorilla/mux"
)

type ValidationRequest struct {
	Profile   string          `json:"profile"`
	JSONLD    json.RawMessage `json:"jsonld"`
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input
	result, err := validator.Validate(req.Profile, string(req.JSONLD), false, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Return the validation result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/validate", validateHandler).Methods("POST")

	port := ":7445"
	log.Printf("Starting ACV validation service on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
