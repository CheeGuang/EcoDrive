package main

import (
	"log"
	"net/http"
	"paymentMicroservice/payment"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Payment endpoints
	router.HandleFunc("/api/v1/payment/real-time-bill", payment.CalculateRealTimeBill).Methods("GET")
	router.HandleFunc("/api/v1/payment/process", payment.ProcessPayment).Methods("POST")
	router.HandleFunc("/api/v1/membership/payment", payment.ProcessMembershipPayment).Methods("POST")


	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5200"}), // Allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type"}),           // Allowed headers
	)(router)

	// Start the server
	log.Println("Payment Microservice is running on port 5200...")
	log.Fatal(http.ListenAndServe(":5200", corsHandler))
}
