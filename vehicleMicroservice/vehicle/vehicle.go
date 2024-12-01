package vehicle

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

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

// Vehicle represents the structure of a vehicle record
type Vehicle struct {
	VehicleID           int     `json:"vehicle_id"`
	Model               string  `json:"model"`
	Location            string  `json:"location"`
	ChargeLevel         *int64  `json:"charge_level,omitempty"`
	CleanlinessStatus   string  `json:"cleanliness_status"`
	RentalPricePerHour  float64 `json:"rental_price_per_hour"`
}

func GetAvailableVehicles(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching available vehicles for specified date range...")

	// Parse query parameters for start_date and end_date
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date query parameters are required", http.StatusBadRequest)
		return
	}

	// Parse ISO 8601 format (used by datetime-local input)
	const iso8601 = "2006-01-02T15:04"
	startDate, err := time.Parse(iso8601, startDateStr)
	if err != nil {
		log.Printf("Invalid start_date format: %v", err)
		http.Error(w, "Invalid start_date format. Use 'YYYY-MM-DDTHH:MM'", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse(iso8601, endDateStr)
	if err != nil {
		log.Printf("Invalid end_date format: %v", err)
		http.Error(w, "Invalid end_date format. Use 'YYYY-MM-DDTHH:MM'", http.StatusBadRequest)
		return
	}

	// Retrieve all vehicles
	rows, err := db.Query("SELECT vehicle_id, model, location, charge_level, cleanliness_status, rental_price_per_hour FROM Vehicles")
	if err != nil {
		log.Printf("Error querying vehicles: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var availableVehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		var chargeLevel sql.NullInt64
		if err := rows.Scan(&vehicle.VehicleID, &vehicle.Model, &vehicle.Location, &chargeLevel, &vehicle.CleanlinessStatus, &vehicle.RentalPricePerHour); err != nil {
			log.Printf("Error scanning vehicle row: %v", err)
			http.Error(w, "Error scanning vehicle row", http.StatusInternalServerError)
			return
		}
		if chargeLevel.Valid {
			vehicle.ChargeLevel = &chargeLevel.Int64
		}

		// Check if the vehicle is booked in the specified date range
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*)
			FROM Bookings
			WHERE vehicle_id = ?
			AND (booking_date < ? AND return_date > ?)`,
			vehicle.VehicleID, endDate, startDate).Scan(&count)
		if err != nil {
			log.Printf("Error checking bookings for vehicle_id %d: %v", vehicle.VehicleID, err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if count == 0 {
			// Vehicle is available
			availableVehicles = append(availableVehicles, vehicle)
			log.Printf("Vehicle available: %+v", vehicle)
		}
	}

	log.Printf("Total available vehicles: %d", len(availableVehicles))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(availableVehicles); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Println("Available vehicles response sent successfully.")
}
// GetVehicleStatus retrieves the status of all vehicles
func GetVehicleStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching vehicle status...")

	rows, err := db.Query("SELECT vehicle_id, model, availability_status, location, charge_level, cleanliness_status, rental_price_per_hour FROM Vehicles")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		var chargeLevel sql.NullInt64
		if err := rows.Scan(&vehicle.VehicleID, &vehicle.Model, &vehicle.Location, &chargeLevel, &vehicle.CleanlinessStatus, &vehicle.RentalPricePerHour); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		if chargeLevel.Valid {
			vehicle.ChargeLevel = &chargeLevel.Int64
		}
		vehicles = append(vehicles, vehicle)
		log.Printf("Fetched vehicle: %+v", vehicle)
	}

	log.Printf("Total vehicles fetched: %d", len(vehicles))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(vehicles); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Println("Vehicle status response sent successfully.")
}
