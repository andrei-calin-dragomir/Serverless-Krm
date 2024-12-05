package server

import (
	"net/http"
	"request_handler/pkg/handlers"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/control-plane", handlers.HandleRequest)
	return router
}
