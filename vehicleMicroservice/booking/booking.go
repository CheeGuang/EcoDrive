package booking

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

// Booking represents the structure of a booking record
type Booking struct {
	BookingID   int     `json:"booking_id"`
	VehicleID   int     `json:"vehicle_id"`
	UserID      int     `json:"user_id"`
	BookingDate string  `json:"booking_date"`
	ReturnDate  string  `json:"return_date"`
	TotalPrice  float64 `json:"total_price"`
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		VehicleID   int     `json:"vehicle_id"`
		UserID      int     `json:"user_id"`
		BookingDate string  `json:"booking_date"`
		ReturnDate  string  `json:"return_date"`
		TotalPrice  float64 `json:"total_price"`
	}

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert the booking into the database
	result, err := db.Exec(`
		INSERT INTO Bookings (vehicle_id, user_id, booking_date, return_date, total_price)
		VALUES (?, ?, ?, ?, ?)`,
		payload.VehicleID, payload.UserID, payload.BookingDate, payload.ReturnDate, payload.TotalPrice)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Retrieve the last inserted booking ID
	bookingID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving booking ID: %v", err)
		http.Error(w, "Failed to retrieve booking ID", http.StatusInternalServerError)
		return
	}

	// Respond with the booking ID
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"booking_id": %d}`, bookingID)))
}


// ModifyBooking allows users to modify an existing booking
func ModifyBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		StartDateTime string  `json:"start_date_time"`
		EndDateTime   string  `json:"end_date_time"`
		TotalPrice    float64 `json:"total_price"`
	}

	// Decode and validate payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if payload.StartDateTime == "" || payload.EndDateTime == "" || payload.TotalPrice <= 0 {
		http.Error(w, "Missing or invalid fields in the input", http.StatusBadRequest)
		return
	}

	// Update the booking in the database
	_, err = db.Exec(`
		UPDATE Bookings 
		SET booking_date = ?, return_date = ?, total_price = ?
		WHERE booking_id = ?`,
		payload.StartDateTime, payload.EndDateTime, payload.TotalPrice, bookingID)
	if err != nil {
		log.Printf("Error updating booking: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Booking updated successfully"))
}


// CancelBooking allows users to cancel an existing booking
func CancelBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM Bookings WHERE booking_id = ?", bookingID)
	if err != nil {
		log.Printf("Error deleting booking: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Booking cancelled successfully"))
}


// GetBooking retrieves details of a specific booking
func GetBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var booking struct {
		BookingID      int     `json:"booking_id"`
		VehicleID      int     `json:"vehicle_id"`
		UserID         int     `json:"user_id"`
		BookingDate    string  `json:"booking_date"`
		ReturnDate     string  `json:"return_date"`
		TotalPrice     float64 `json:"total_price"`
		Model          string  `json:"model"`
		Location       string  `json:"location"`
		ChargeLevel    int     `json:"charge_level"`
	}

	err = db.QueryRow(`
		SELECT 
			b.booking_id, b.vehicle_id, b.user_id, 
			b.booking_date, b.return_date, b.total_price,
			v.model, v.location, v.charge_level
		FROM Bookings b
		JOIN Vehicles v ON b.vehicle_id = v.vehicle_id
		WHERE b.booking_id = ?`, bookingID).
		Scan(
			&booking.BookingID,
			&booking.VehicleID,
			&booking.UserID,
			&booking.BookingDate,
			&booking.ReturnDate,
			&booking.TotalPrice,
			&booking.Model,
			&booking.Location,
			&booking.ChargeLevel,
		)

	if err == sql.ErrNoRows {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// GetBookingsByUserID retrieves all bookings for a specific user
func GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBookingsByUserID: Start processing request") // Debug: Start of function

	params := mux.Vars(r)
	log.Printf("GetBookingsByUserID: Extracted params: %v\n", params) // Debug: Log params

	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		log.Printf("GetBookingsByUserID: Invalid user ID: %v\n", params["user_id"]) // Debug: Invalid ID
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("GetBookingsByUserID: Converted user ID: %d\n", userID) // Debug: Valid ID

	rows, err := db.Query(`
		SELECT 
			b.booking_id, b.vehicle_id, b.user_id, 
			b.booking_date, b.return_date, b.total_price,
			v.model, v.location, v.charge_level, v.rental_price_per_hour
		FROM Bookings b
		JOIN Vehicles v ON b.vehicle_id = v.vehicle_id
		WHERE b.user_id = ?`, userID)
	if err != nil {
		log.Printf("GetBookingsByUserID: Error retrieving bookings: %v\n", err) // Debug: Query error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	log.Println("GetBookingsByUserID: Query executed successfully") // Debug: Query success

	var bookings []struct {
		BookingID         int     `json:"booking_id"`
		VehicleID         int     `json:"vehicle_id"`
		UserID            int     `json:"user_id"`
		BookingDate       string  `json:"booking_date"`
		ReturnDate        string  `json:"return_date"`
		TotalPrice        float64 `json:"total_price"`
		Model             string  `json:"model"`
		Location          string  `json:"location"`
		ChargeLevel       int     `json:"charge_level"`
		RentalPricePerHour float64 `json:"rental_price_per_hour"`
	}

	for rows.Next() {
		var booking struct {
			BookingID         int     `json:"booking_id"`
			VehicleID         int     `json:"vehicle_id"`
			UserID            int     `json:"user_id"`
			BookingDate       string  `json:"booking_date"`
			ReturnDate        string  `json:"return_date"`
			TotalPrice        float64 `json:"total_price"`
			Model             string  `json:"model"`
			Location          string  `json:"location"`
			ChargeLevel       int     `json:"charge_level"`
			RentalPricePerHour float64 `json:"rental_price_per_hour"`
		}
		if err := rows.Scan(
			&booking.BookingID,
			&booking.VehicleID,
			&booking.UserID,
			&booking.BookingDate,
			&booking.ReturnDate,
			&booking.TotalPrice,
			&booking.Model,
			&booking.Location,
			&booking.ChargeLevel,
			&booking.RentalPricePerHour,
		); err != nil {
			log.Printf("GetBookingsByUserID: Error scanning row: %v\n", err) // Debug: Scan error
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		log.Printf("GetBookingsByUserID: Retrieved booking: %+v\n", booking) // Debug: Retrieved booking
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetBookingsByUserID: Row iteration error: %v\n", err) // Debug: Row iteration error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("GetBookingsByUserID: Total bookings retrieved: %d\n", len(bookings)) // Debug: Number of bookings
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("GetBookingsByUserID: Error encoding response: %v\n", err) // Debug: Encode error
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Println("GetBookingsByUserID: Response sent successfully") // Debug: End of function
}

// GetBookingsByVehicleID retrieves booking dates for a specific vehicle
func GetBookingsByVehicleID(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBookingsByVehicleID: Start processing request") // Debug: Start of function

	params := mux.Vars(r)
	log.Printf("GetBookingsByVehicleID: Extracted params: %v\n", params) // Debug: Log params

	vehicleID, err := strconv.Atoi(params["vehicle_id"])
	if err != nil {
		log.Printf("GetBookingsByVehicleID: Invalid vehicle ID: %v\n", params["vehicle_id"]) // Debug: Invalid ID
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}
	log.Printf("GetBookingsByVehicleID: Converted vehicle ID: %d\n", vehicleID) // Debug: Valid ID

	rows, err := db.Query(`
		SELECT 
			b.booking_date, b.return_date
		FROM Bookings b
		WHERE b.vehicle_id = ?`, vehicleID)
	if err != nil {
		log.Printf("GetBookingsByVehicleID: Error retrieving bookings by vehicle ID: %v\n", err) // Debug: Query error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	log.Println("GetBookingsByVehicleID: Query executed successfully") // Debug: Query success

	var bookings []struct {
		BookingDate string `json:"booking_date"`
		ReturnDate  string `json:"return_date"`
	}

	for rows.Next() {
		var booking struct {
			BookingDate string `json:"booking_date"`
			ReturnDate  string `json:"return_date"`
		}
		if err := rows.Scan(&booking.BookingDate, &booking.ReturnDate); err != nil {
			log.Printf("GetBookingsByVehicleID: Error scanning row: %v\n", err) // Debug: Scan error
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		log.Printf("GetBookingsByVehicleID: Retrieved booking: %+v\n", booking) // Debug: Retrieved booking
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetBookingsByVehicleID: Row iteration error: %v\n", err) // Debug: Row error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("GetBookingsByVehicleID: Total bookings retrieved: %d\n", len(bookings)) // Debug: Number of bookings
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("GetBookingsByVehicleID: Error encoding response: %v\n", err) // Debug: Encode error
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	log.Println("GetBookingsByVehicleID: Response sent successfully") // Debug: End of function
}