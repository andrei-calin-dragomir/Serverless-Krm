package server

import (
	"net/http"
	"serverless_functions/pkg/handlers"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/control-plane", handlers.HandleRequest)
	return router
}
