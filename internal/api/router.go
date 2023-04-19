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
	router.HandleFunc("/api/v1/health", apiController.HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/hook", apiController.ReceiveMessage).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/hook", apiController.VerifyToken).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/{number}/document", apiController.UploadDocument).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/{number}/document", apiController.GetDocument).Methods(http.MethodGet)
	// Add rate limiting middleware to all endpoints
	return router
}
