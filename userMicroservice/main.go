package main

import (
	"log"
	"net/http"
	"userMicroservice/membership"
	"userMicroservice/profile"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the router
	router := mux.NewRouter()

	// Membership endpoints
	router.HandleFunc("/api/v1/user/membership/status", membership.GetMembershipStatus).Methods("GET")
	router.HandleFunc("/api/v1/user/membership/update", membership.UpdateMembershipTier).Methods("PUT")

	// Profile management endpoints
	router.HandleFunc("/api/v1/user/create", profile.CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/profile", profile.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/v1/user/profile/update", profile.UpdateUserProfile).Methods("PUT")
	router.HandleFunc("/api/v1/user/rental-history", profile.GetRentalHistory).Methods("GET")

	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5100"}), // Update for allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}), // Update for allowed HTTP methods
		handlers.AllowedHeaders([]string{"Content-Type"}), // Update for allowed headers
	)(router)

	// Start the server
	log.Println("User Microservice is running on port 5100...")
	log.Fatal(http.ListenAndServe(":5100", corsHandler))
}
