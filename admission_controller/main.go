package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
)

var RECONCILER_URL string

func init() {
	RECONCILER_URL = os.Getenv("RECONCILER_URL")
}

func main() {
	// Set up a simple HTTP server
	http.HandleFunc("/validate", handleValidate)

	// Set the server to listen on port 443 for incoming requests
	port := "443"

	fmt.Printf("Starting server on :%s...\n", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

// Handle validation requests from the Kubernetes API server
func handleValidate(w http.ResponseWriter, r *http.Request) {

	var admissionReview admissionv1.AdmissionReview

	// Decode the request body into the AdmissionReview object
	if err := json.NewDecoder(r.Body).Decode(&admissionReview); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode admission review: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the response AdmissionReview
	admissionResponse := admissionv1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: true, // Allow by default
	}

	resp := admissionv1.AdmissionReview{
		Response: &admissionResponse,
	}

	// Respond back to the API server with the AdmissionReview response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}

	// pod := admissionReview.Request.Object.Raw
	// fmt.Printf("Received Pod Update:\n%s\n", pod)
	// Extract the namespace and name from the AdmissionReview
	namespace := admissionReview.Request.Namespace
	name := admissionReview.Request.Name

	jsonBytes, err := json.Marshal(struct {
		Kind      string `json:"kind"`
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}{
		Kind:      "pod",
		Namespace: namespace,
		Name:      name,
	})

	if err != nil {
		_ = err
	}

	response, err := PostJSON(RECONCILER_URL, jsonBytes)

	if err != nil {
		slog.Error("Error sending pod update to reconciler",
			slog.Any("error", err),
			slog.String("url", RECONCILER_URL),
		)

	}
	slog.Info("Received response from reconciler",
		slog.Any("response", response),
	)

	// Here you can implement your custom logic to allow or reject the Pod update
	// If you want to reject the update, you can set Allowed to false and provide a rejection message
	// admissionResponse.Allowed = false
	// admissionResponse.Result = &metav1.Status{
	// 	Code:    403,
	// 	Message: "Custom rejection message",
	// }

}

func PostJSON(url string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
