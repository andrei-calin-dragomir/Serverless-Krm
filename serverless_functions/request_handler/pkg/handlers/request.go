package handlers

import (
	"encoding/json"         // For JSON parsing
	"fmt"                   // For formatted output
	"log"                   // For logging
	"net/http"              // For HTTP handling

	appsv1 "k8s.io/api/apps/v1"            // Kubernetes Deployment struct
	batchv1 "k8s.io/api/batch/v1" // Job struct
)

// Supported resource types
type ResourceType string;
const (
	Deployment ResourceType = "Deployment"
	Job ResourceType = "Job"
)

type RequestPayload struct {
	Resource string      `json:"resource"` // Resource type, e.g., "Deployment"
	Action   string      `json:"action"`   // Action, e.g., "create"
	Object   interface{} `json:"object"`   // Dynamic object field
}
	
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received request: %+v\n", payload)

	// Switch based on resource type
	switch ResourceType(payload.Resource) {
	case Deployment:
		var deployment appsv1.Deployment
		if err := decodeObject(payload.Object, &deployment); err != nil {
			http.Error(w, "Invalid Deployment object", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed Deployment: %+v\n", deployment)

	case Job:
		var job batchv1.Job
		if err := decodeObject(payload.Object, &job); err != nil {
			http.Error(w, "Invalid Job object", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed Job: %+v\n", job)

	default:
		http.Error(w, "Unsupported resource type", http.StatusNotImplemented)
		return
	}

	// Acknowledge the request
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Request for %s action on %s received", payload.Action, payload.Resource)
}

// Utility to decode dynamic objects into specific Kubernetes structs
func decodeObject(object interface{}, out interface{}) error {
	bytes, err := json.Marshal(object) // Marshal dynamic object to bytes
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out) // Unmarshal into target struct
}
