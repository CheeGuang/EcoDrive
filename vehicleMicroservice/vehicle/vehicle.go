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

// UpdateVehicleDetails updates cleanliness status, charge level, and location of a vehicle,
// and optionally updates the booking status to "active".
func UpdateVehicleDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting UpdateVehicleDetails function...")

	// Check the HTTP method
	if r.Method != http.MethodPut {
		log.Println("Invalid HTTP method. Only PUT is allowed.")
		http.Error(w, "Invalid request method. Use PUT.", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body
	var updateRequest struct {
		VehicleID         int     `json:"vehicle_id"`
		CleanlinessStatus string  `json:"cleanliness_status,omitempty"`
		ChargeLevel       *int64  `json:"charge_level,omitempty"`
		Location          string  `json:"location,omitempty"`
		BookingID         *int    `json:"booking_id,omitempty"` // Optional booking ID for updating status
	}
	log.Println("Parsing request body...")
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		log.Printf("Error decoding JSON payload: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Parsed request: %+v", updateRequest)

	// Validate required fields
	if updateRequest.VehicleID == 0 {
		log.Println("Vehicle ID is missing in the request payload.")
		http.Error(w, "Vehicle ID is required", http.StatusBadRequest)
		return
	}

	// Prepare the SQL query dynamically based on provided fields
	query := "UPDATE Vehicles SET "
	var args []interface{}
	if updateRequest.CleanlinessStatus != "" {
		log.Printf("Adding cleanliness_status to update query: %s", updateRequest.CleanlinessStatus)
		query += "cleanliness_status = ?, "
		args = append(args, updateRequest.CleanlinessStatus)
	}
	if updateRequest.ChargeLevel != nil {
		log.Printf("Adding charge_level to update query: %d", *updateRequest.ChargeLevel)
		query += "charge_level = ?, "
		args = append(args, *updateRequest.ChargeLevel)
	}
	if updateRequest.Location != "" {
		log.Printf("Adding location to update query: %s", updateRequest.Location)
		query += "location = ?, "
		args = append(args, updateRequest.Location)
	}
	// Remove trailing comma and space
	query = query[:len(query)-2]
	query += " WHERE vehicle_id = ?"
	args = append(args, updateRequest.VehicleID)

	log.Printf("Generated query for vehicle update: %s", query)
	log.Printf("Query arguments: %+v", args)

	// Execute the query to update the vehicle
	log.Println("Executing vehicle update query...")
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("Error executing vehicle update query: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		http.Error(w, "Error fetching update result", http.StatusInternalServerError)
		return
	}
	log.Printf("Rows affected by vehicle update: %d", rowsAffected)
	if rowsAffected == 0 {
		log.Printf("No rows updated for vehicle_id: %d", updateRequest.VehicleID)
		http.Error(w, "Vehicle not found or no changes made", http.StatusNotFound)
		return
	}

	// If a BookingID is provided, update the booking status to "active"
	if updateRequest.BookingID != nil {
		log.Printf("Updating booking status to 'active' for booking_id: %d", *updateRequest.BookingID)
		_, err := db.Exec(`
			UPDATE Bookings 
			SET status = 'active' 
			WHERE booking_id = ?`, *updateRequest.BookingID)
		if err != nil {
			log.Printf("Error updating booking status: %v", err)
			http.Error(w, "Database error while updating booking status", http.StatusInternalServerError)
			return
		}
		log.Printf("Booking status successfully updated to 'active' for booking_id: %d", *updateRequest.BookingID)
	} else {
		log.Println("No booking_id provided. Skipping booking status update.")
	}

	// Send success response
	log.Println("Sending success response...")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vehicle details and booking status updated successfully",
	})

	log.Println("UpdateVehicleDetails function completed.")
}