package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdConfig holds the configuration for connecting to the etcd cluster
type EtcdConfig struct {
	Endpoints   []string
	DialTimeout time.Duration
}

// AuthorizationRequest represents the incoming request structure
type AuthorizationRequest struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// AuthorizationResponse represents the structure of the response
type AuthorizationResponse struct {
	Authorized bool   `json:"authorized"`
	Message    string `json:"message"`
}

// getEtcdConfigFromEnv retrieves etcd configuration from environment variables
func getEtcdConfigFromEnv() (EtcdConfig, error) {
	endpoints := os.Getenv("ETCD_ENDPOINTS")
	if endpoints == "" {
		return EtcdConfig{}, fmt.Errorf("ETCD_ENDPOINTS environment variable not set")
	}

	// Split endpoints by comma
	endpointList := strings.Split(endpoints, ",")

	return EtcdConfig{
		Endpoints:   endpointList,
		DialTimeout: 5 * time.Second,
	}, nil
}

// AuthorizeHandler handles the authorization of a request
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var authReq AuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input
	if authReq.Token == "" || authReq.User == "" {
		http.Error(w, "Token and User fields are required", http.StatusBadRequest)
		return
	}

	// Load etcd configuration
	etcdConfig, err := getEtcdConfigFromEnv()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load etcd configuration: %v", err), http.StatusInternalServerError)
		return
	}

	// Initialize etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConfig.Endpoints,
		DialTimeout: etcdConfig.DialTimeout,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to etcd: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Check authorization data in etcd
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Assume keys in etcd are structured as "auth/<user>"
	key := fmt.Sprintf("auth/%s", authReq.User)
	resp, err := client.Get(ctx, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve data from etcd: %v", err), http.StatusInternalServerError)
		return
	}

	if len(resp.Kvs) == 0 {
		// User not found in etcd
		json.NewEncoder(w).Encode(AuthorizationResponse{
			Authorized: false,
			Message:    "User not found",
		})
		return
	}

	// Extract and validate token
	storedToken := string(resp.Kvs[0].Value)
	if storedToken != authReq.Token {
		json.NewEncoder(w).Encode(AuthorizationResponse{
			Authorized: false,
			Message:    "Invalid token",
		})
		return
	}

	// If valid
	json.NewEncoder(w).Encode(AuthorizationResponse{
		Authorized: true,
		Message:    "Authorization successful",
	})
}