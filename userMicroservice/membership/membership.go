package membership

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

    // Fetch DB connection string from environment variables
    dbConnection := os.Getenv("DB_CONNECTION")
    if dbConnection == "" {
        log.Fatalf("DB_CONNECTION environment variable is not set")
    }
    log.Printf("DB_CONNECTION (membership package): %s", dbConnection) // Debugging

    // Initialize the database connection
    db, err = sql.Open("mysql", dbConnection)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    // Test the database connection
    err = db.Ping()
    if err != nil {
        log.Fatalf("Database connection test failed: %v", err)
    }
    log.Println("Database connection (membership package) successful.")
}

// GetMembershipStatus retrieves the user's membership status
func GetMembershipStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var membershipTier string
	err := db.QueryRow("SELECT membership_level FROM User WHERE user_id = ?", userID).Scan(&membershipTier)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"membership_level": membershipTier})
}

// UpdateMembershipTier updates the membership level of a user
func UpdateMembershipTier(w http.ResponseWriter, r *http.Request) {
	dbConnection := os.Getenv("DB_CONNECTION")
if dbConnection == "" {
    log.Fatalf("DB_CONNECTION environment variable is not set")
}
log.Printf("DB_CONNECTION: %s", dbConnection) // Debug line


	log.Println("Received request to update membership tier")

	var payload struct {
		UserID         int    `json:"user_id"`
		MembershipTier string `json:"membership_tier"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Invalid input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded payload: user_id=%d, membership_tier=%s", payload.UserID, payload.MembershipTier)

	// Execute the SQL update
	query := "UPDATE User SET membership_level = ? WHERE user_id = ?"
	log.Printf("Executing query: %s with values (%s, %d)", query, payload.MembershipTier, payload.UserID)
	result, err := db.Exec(query, payload.MembershipTier, payload.UserID)
	if err != nil {
		log.Printf("Error updating membership level: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("Rows affected: %d", rowsAffected)

	if rowsAffected == 0 {
		log.Printf("No user found with user_id=%d", payload.UserID)
		http.Error(w, "No user found with the given ID", http.StatusNotFound)
		return
	}

	// Return success response
	log.Printf("Membership level successfully updated for user_id=%d", payload.UserID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Membership level updated successfully"}`))
}