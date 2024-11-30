package membership

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// Initialize database connection
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
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
	var payload struct {
		UserID         int    `json:"user_id"`
		MembershipTier string `json:"membership_tier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE User SET membership_level = ? WHERE user_id = ?", payload.MembershipTier, payload.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Membership level updated successfully"))
}
