package main

import (
	"log"
	"net/http"
	"os"
	"paymentMicroservice/payment"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Authentication middleware to validate JWT
func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Get the JWT secret from the environment variable
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			log.Println("JWT_SECRET is not set in the environment")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid JWT token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	// JWT Authentication Logic
	authenticated := router.NewRoute().Subrouter()
	authenticated.Use(authenticateMiddleware)

	// Payment endpoints
	authenticated.HandleFunc("/api/v1/payment/real-time-bill", payment.CalculateRealTimeBill).Methods("GET")
	authenticated.HandleFunc("/api/v1/payment/process", payment.ProcessPayment).Methods("POST")
	authenticated.HandleFunc("/api/v1/membership/payment", payment.ProcessMembershipPayment).Methods("POST")

	// Add CORS support
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5200"}), // Allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}), // Allowed methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Include Authorization header
	)(router)

	// Start the server
	log.Println("Payment Microservice is running on port 5200...")
	log.Fatal(http.ListenAndServe(":5200", corsHandler))
}