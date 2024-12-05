package server

import (
	"net/http"
	"authorization/pkg/handlers"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/authorize", handlers.AuthorizeHandler)
	return router
}
