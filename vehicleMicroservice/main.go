package main

import (
	"log"
	"net/http"
	"vehicleMicroservice/booking"
	"vehicleMicroservice/vehicle"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the router
	router := mux.NewRouter()

	// Vehicle endpoints
	router.HandleFunc("/api/v1/vehicle/availability", vehicle.GetAvailableVehicles).Methods("GET")
	router.HandleFunc("/api/v1/vehicle/status", vehicle.GetVehicleStatus).Methods("GET")

	// Booking endpoints
	router.HandleFunc("/api/v1/vehicle/booking", booking.CreateBooking).Methods("POST")
	router.HandleFunc("/api/v1/vehicle/booking/{id}", booking.GetBooking).Methods("GET")
	router.HandleFunc("/api/v1/vehicle/booking/{id}", booking.ModifyBooking).Methods("PUT")
	router.HandleFunc("/api/v1/vehicle/booking/{id}", booking.CancelBooking).Methods("DELETE")
	router.HandleFunc("/api/v1/vehicle/booking/user/{user_id}", booking.GetBookingsByUserID).Methods("GET")
	router.HandleFunc("/api/v1/vehicle/booking/vehicle/{vehicle_id}", booking.GetBookingsByVehicleID).Methods("GET")



	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5150"}), // Allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type"}), // Allowed headers
	)(router)

	// Start the server
	log.Println("Vehicle Microservice is running on port 5150...")
	log.Fatal(http.ListenAndServe(":5150", corsHandler))
}
