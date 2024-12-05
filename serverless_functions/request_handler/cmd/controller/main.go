package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	batchv1 "k8s.io/api/batch/v1" // Import Kubernetes Job struct
)

// RequestPayload is the expected JSON payload for incoming requests
type RequestPayload struct {
	Resource string      `json:"resource"`
	Action   string      `json:"action"`
	Object   interface{} `json:"object"`
}

// decodeObject converts a generic interface{} into the specified Kubernetes object
func decodeObject(object interface{}, out interface{}) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming request
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Process resource type
	switch payload.Resource {
	case "Job":
		var job batchv1.Job
		if err := decodeObject(payload.Object, &job); err != nil {
			http.Error(w, "Invalid Job object", http.StatusBadRequest)
			return
		}
		log.Printf("Received Job: %+v\n", job)

	default:
		http.Error(w, "Unsupported resource type", http.StatusNotImplemented)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Request accepted and processed")
}

func main() {
	http.HandleFunc("/", handleRequest)

	port := "8080"
	log.Printf("Starting serverless controller on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
