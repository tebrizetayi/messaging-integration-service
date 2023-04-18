package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const ()

// NewAPI returns a new API router
// The router is configured with the API controller
// and the rate limiting middleware
func NewAPI(apiController Controller) http.Handler {
	router := mux.NewRouter()

	// Add API endpoints
	router.HandleFunc("/api/v1/hook", apiController.ReceiveMessage).Methods(http.MethodPost)

	// Add rate limiting middleware to all endpoints
	return router
}
