package profile

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	dbConnection := os.Getenv("DB_CONNECTION")
	if dbConnection == "" {
		log.Fatalf("DB_CONNECTION environment variable is not set")
	}

	log.Println("Initializing database connection...")
	db, err = sql.Open("mysql", dbConnection)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}
	log.Println("Database connection successful.")
}

// CreateUserRequest represents the structure of the request to create a new user
type CreateUserRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	ContactNumber string `json:"contact_number"`
	Address       string `json:"address"`
	Password      string `json:"password"` // Assuming password is hashed before sending
}

// CreateUser handles the creation of a new user record
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert the new user record with default membership level as 'Basic'
	_, err = db.Exec(`
		INSERT INTO User (name, email, contact_number, address, password, membership_level)
		VALUES (?, ?, ?, ?, ?, 'Basic')`,
		req.Name, req.Email, req.ContactNumber, req.Address, req.Password,
	)

	if err != nil {
		log.Printf("Error inserting user record: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// GetUserProfile retrieves a user's profile information
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var profile struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Contact    string `json:"contact_number"`
		Address    string `json:"address"`
		Membership string `json:"membership_level"`
	}

	err := db.QueryRow("SELECT name, email, contact_number, address, membership_level FROM User WHERE user_id = ?", userID).Scan(&profile.Name, &profile.Email, &profile.Contact, &profile.Address, &profile.Membership)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateUserProfile updates a user's personal details
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var profile struct {
		UserID  int    `json:"user_id"`
		Name    string `json:"name"`
		Contact string `json:"contact_number"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE User SET name = ?, contact_number = ?, address = ? WHERE user_id = ?", profile.Name, profile.Contact, profile.Address, profile.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile updated successfully"))
}

// GetRentalHistory retrieves a user's rental history
func GetRentalHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	rows, err := db.Query("SELECT vehicle_id, rental_price_per_hour, created_at FROM Rentals WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []struct {
		VehicleID        int     `json:"vehicle_id"`
		RentalPrice      float64 `json:"rental_price_per_hour"`
		RentalDate       string  `json:"rental_date"`
	}
	for rows.Next() {
		var record struct {
			VehicleID   int     `json:"vehicle_id"`
			RentalPrice float64 `json:"rental_price_per_hour"`
			RentalDate  string  `json:"rental_date"`
		}
		if err := rows.Scan(&record.VehicleID, &record.RentalPrice, &record.RentalDate); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		history = append(history, record)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
