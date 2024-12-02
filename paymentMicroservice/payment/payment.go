package payment

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jung-kurt/gofpdf"
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

// TierBasedPricing calculates pricing based on membership level and rental duration
func TierBasedPricing(userID int, durationHours int, pricePerHour float64) (float64, float64) {
	var membershipLevel string
	var discountPercentage float64

	// Fetch membership level
	err := db.QueryRow("SELECT membership_level FROM ecoDrive_user_db.User WHERE user_id = ?", userID).Scan(&membershipLevel)
	if err != nil {
		log.Printf("Error fetching membership level: %v", err)
		return pricePerHour * float64(durationHours), 0.00
	}

	// Fetch discount percentage
	err = db.QueryRow("SELECT discount_percentage FROM Discounts WHERE membership_level = ?", membershipLevel).Scan(&discountPercentage)
	if err != nil {
		log.Printf("Error fetching discount percentage: %v", err)
		return pricePerHour * float64(durationHours), 0.00
	}

	// Calculate total price and discount
	totalPrice := pricePerHour * float64(durationHours)
	discount := totalPrice * (discountPercentage / 100)
	finalPrice := totalPrice - discount
	return finalPrice, discount
}

// CalculateRealTimeBill handles real-time billing calculation
func CalculateRealTimeBill(w http.ResponseWriter, r *http.Request) {
	membershipLevel := r.URL.Query().Get("membership_level")
	durationHours, _ := strconv.Atoi(r.URL.Query().Get("duration_hours"))
	pricePerHour, _ := strconv.ParseFloat(r.URL.Query().Get("price_per_hour"), 64)

	var discountPercentage float64

	// Fetch discount percentage
	err := db.QueryRow("SELECT discount_percentage FROM Discounts WHERE membership_level = ?", membershipLevel).Scan(&discountPercentage)
	if err != nil {
		log.Printf("Error fetching discount percentage: %v", err)
		http.Error(w, "Invalid membership level", http.StatusBadRequest)
		return
	}

	// Calculate total price and discount
	totalPrice := pricePerHour * float64(durationHours)
	discount := totalPrice * (discountPercentage / 100)
	finalPrice := totalPrice - discount

	response := map[string]interface{}{
		"final_price":  finalPrice,
		"discount":     discount,
		"total_price":  totalPrice,
		"membership":   membershipLevel,
		"duration":     durationHours,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var payment struct {
		UserID        int     `json:"user_id"`
		VehicleID     string  `json:"vehicle_id"`
		StartDate     string  `json:"start_date"`
		EndDate       string  `json:"end_date"`
		PaymentMethod string  `json:"payment_method"`
		PricePerHour  string  `json:"price_per_hour"`
		RentalDuration string `json:"rental_duration"`
		TotalPrice    string  `json:"total_price"`
		Email         string  `json:"email"`
	}

	// Decode incoming JSON request
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		log.Printf("Error decoding payment request: %v", err)
		http.Error(w, "Invalid payment request", http.StatusBadRequest)
		return
	}
	log.Printf("Decoded payment payload: %+v", payment)

	// Parse required fields to correct types
	vehicleID, err := strconv.Atoi(payment.VehicleID)
	if err != nil {
		log.Printf("Error parsing vehicle_id: %v", err)
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	totalPrice, err := strconv.ParseFloat(payment.TotalPrice, 64)
	if err != nil {
		log.Printf("Error parsing total price: %v", err)
		http.Error(w, "Invalid total price", http.StatusBadRequest)
		return
	}

	// Convert dates to expected format (if needed)
	startDate, err := time.Parse("2006-01-02T15:04", payment.StartDate)
	if err != nil {
		log.Printf("Error parsing start_date: %v", err)
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02T15:04", payment.EndDate)
	if err != nil {
		log.Printf("Error parsing end_date: %v", err)
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	// Notify booking service
	apiURL := "http://localhost:5150/api/v1/vehicle/booking"
	bookingPayload := map[string]interface{}{
		"vehicle_id":   vehicleID,
		"user_id":      payment.UserID,
		"booking_date": startDate.Format("2006-01-02 15:04:05"),
		"return_date":  endDate.Format("2006-01-02 15:04:05"),
		"total_price":  totalPrice,
	}

	// Log the booking payload for debugging
	log.Printf("Booking payload: %+v", bookingPayload)

	jsonPayload, _ := json.Marshal(bookingPayload)
	log.Printf("Serialized JSON payload: %s", string(jsonPayload))

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling booking API: %v", err)
		http.Error(w, "Failed to notify booking service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Booking API returned non-OK status: %d, Response: %s", resp.StatusCode, string(body))
		http.Error(w, "Failed to notify booking service", http.StatusInternalServerError)
		return
	}

	// Parse the booking ID from the API response
	var bookingResponse struct {
		BookingID int `json:"booking_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		log.Printf("Error decoding booking API response: %v", err)
		http.Error(w, "Failed to process booking response", http.StatusInternalServerError)
		return
	}

	log.Printf("Received booking ID from API: %d", bookingResponse.BookingID)

	// Insert payment details into the database
	result, err := db.Exec(`
			INSERT INTO BookingPayment (user_id, booking_id, amount, payment_method, payment_status, discount, final_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
		payment.UserID, bookingResponse.BookingID, totalPrice, payment.PaymentMethod, "Completed", 0.00, totalPrice)
	if err != nil {
		log.Printf("Error storing payment details: %v", err)
		http.Error(w, "Failed to store payment details", http.StatusInternalServerError)
		return
	}

	// Retrieve the last inserted payment ID
	paymentID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving payment ID: %v", err)
		http.Error(w, "Failed to retrieve payment ID", http.StatusInternalServerError)
		return
	}

	log.Printf("Payment processed successfully for user_id: %d, booking_id: %d, payment_id: %d", payment.UserID, bookingResponse.BookingID, paymentID)

	// Generate and send invoice
	err = generateInvoiceAndSendEmail(
		bookingResponse.BookingID,
		int(paymentID),
		payment.UserID,
		totalPrice,
		payment.PaymentMethod,
		payment.Email,
		startDate,
		endDate,
	)
	if err != nil {
		log.Printf("Error sending invoice email: %v", err)
		http.Error(w, "Failed to send invoice email", http.StatusInternalServerError)
		return
	}

	// Respond with a JSON object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Payment processed successfully",
		"booking_id": bookingResponse.BookingID,
		"payment_id": paymentID,
	})
}


// generateInvoice generates an invoice PDF and returns it as a byte slice
func generateInvoice(bookingID int, paymentID int, userID int, totalPrice float64, paymentMethod string, startDate, endDate time.Time) ([]byte, error) {
    // Create a new PDF document
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)

    // Header
    pdf.SetFillColor(25, 135, 84) // Green colour scheme
    pdf.SetTextColor(255, 255, 255)
    pdf.CellFormat(0, 10, "EcoDrive Invoice", "1", 1, "C", true, 0, "")

    // Add content
    pdf.SetTextColor(0, 0, 0)
    pdf.SetFont("Arial", "", 12)

    // Booking and payment details
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Booking ID: %d", bookingID))
    pdf.Ln(6)
    pdf.Cell(40, 10, fmt.Sprintf("Payment ID: %d", paymentID))
    pdf.Ln(6)
    pdf.Cell(40, 10, fmt.Sprintf("User ID: %d", userID))
    pdf.Ln(6)
    pdf.Cell(40, 10, fmt.Sprintf("Total Price: $%.2f", totalPrice))
    pdf.Ln(6)
    pdf.Cell(40, 10, fmt.Sprintf("Payment Method: %s", paymentMethod))
    pdf.Ln(6)
    pdf.Cell(40, 10, fmt.Sprintf("Start Date: %s", startDate.Format("2006-01-02 15:04:05")))
    pdf.Ln(6)
    pdf.Cell(40, 10, fmt.Sprintf("End Date: %s", endDate.Format("2006-01-02 15:04:05")))

    // Footer
    pdf.Ln(20)
    pdf.SetFont("Arial", "I", 10)
    pdf.SetTextColor(128, 128, 128)
    pdf.CellFormat(0, 10, "Thank you for choosing EcoDrive. Drive safe!", "", 1, "C", false, 0, "")

    // Write PDF to memory buffer
    buf := new(bytes.Buffer)
    err := pdf.Output(buf)
    if err != nil {
        log.Printf("Error generating PDF: %v", err)
        return nil, err
    }

    return buf.Bytes(), nil
}



// sendEmailWithAttachment sends an email with a PDF attachment
func sendEmailWithAttachment(to, subject, body string, fileName string, fileBytes []byte) error {
    // SMTP configuration
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    smtpUser := os.Getenv("SMTP_USER")
    smtpPassword := os.Getenv("SMTP_PASSWORD")

    from := "EcoDrive <" + smtpUser + ">"

    // Create the email with the attachment
    boundary := "EcoDriveBoundary"
    message := bytes.NewBuffer(nil)
    message.WriteString(fmt.Sprintf("From: %s\r\n", from))
    message.WriteString(fmt.Sprintf("To: %s\r\n", to))
    message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
    message.WriteString("MIME-Version: 1.0\r\n")
    message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
    message.WriteString("\r\n--" + boundary + "\r\n")
    message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
    message.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
    message.WriteString(body)
    message.WriteString("\r\n--" + boundary + "\r\n")
    message.WriteString("Content-Type: application/pdf\r\n")
    message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
    message.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	message.WriteString(encodeToBase64(fileBytes))
    message.WriteString("\r\n--" + boundary + "--\r\n")

    // Send email
    auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, message.Bytes())
    if err != nil {
        log.Printf("Error sending email with attachment: %v", err)
        return err
    }

    log.Println("Email sent successfully.")
    return nil
}

// generateInvoiceAndSendEmail generates an invoice and sends it as an email attachment
func generateInvoiceAndSendEmail(bookingID int, paymentID int, userID int, totalPrice float64, paymentMethod, userEmail string, startDate, endDate time.Time) error {
    // Generate the invoice in memory
    fileBytes, err := generateInvoice(bookingID, paymentID, userID, totalPrice, paymentMethod, startDate, endDate)
    if err != nil {
        return err
    }

    // Email content
    subject := "Your EcoDrive Invoice"
    body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Invoice</title>
        </head>
        <body>
            <p>Dear User,</p>
            <p>Thank you for using EcoDrive! Attached is your invoice for the recent transaction.</p>
            <p>Details:</p>
            <ul>
                <li>Booking ID: %d</li>
                <li>Payment ID: %d</li>
                <li>Total Price: $%.2f</li>
            </ul>
            <p>We hope you had a pleasant experience!</p>
            <p>Best regards,<br>The EcoDrive Team</p>
        </body>
        </html>
    `, bookingID, paymentID, totalPrice)

    // Send the email with the invoice attached
    fileName := fmt.Sprintf("Invoice_%d.pdf", paymentID)
    return sendEmailWithAttachment(userEmail, subject, body, fileName, fileBytes)
}

// encodeToBase64 encodes bytes to a base64 string
func encodeToBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// AddMembershipPayment adds a payment entry to the MembershipPayment table
func AddMembershipPayment(userID int, membershipLevel string, amount float64, paymentMethod string, startDate, endDate time.Time) (int64, error) {
    result, err := db.Exec(`
        INSERT INTO MembershipPayment (user_id, membership_level, amount, payment_method, payment_status, start_date, end_date)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
        userID, membershipLevel, amount, paymentMethod, "Completed", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    if err != nil {
        log.Printf("Error inserting membership payment: %v", err)
        return 0, err
    }

    // Retrieve the last inserted payment ID
    membershipPaymentID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error retrieving membership payment ID: %v", err)
        return 0, err
    }

    log.Printf("Membership payment successfully added. Payment ID: %d", membershipPaymentID)
    return membershipPaymentID, nil
}

func ProcessMembershipPayment(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to process membership payment")

	var payment struct {
		UserID          int     `json:"user_id"`
		MembershipLevel string  `json:"membership_level"`
		Amount          float64 `json:"amount"`
		PaymentMethod   string  `json:"payment_method"`
		StartDate       string  `json:"start_date"`
		EndDate         string  `json:"end_date"`
		Email           string  `json:"email"`
	}

	// Decode incoming JSON request
	log.Println("Decoding incoming JSON request")
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		log.Printf("[ERROR] Decoding payment request: %v", err)
		http.Error(w, "Invalid membership payment request", http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG] Decoded payment request: %+v", payment)

	// Parse dates to expected format
	log.Println("[DEBUG] Parsing start_date and end_date")
	startDate, err := time.Parse("2006-01-02", payment.StartDate)
	if err != nil {
		log.Printf("[ERROR] Parsing start_date: %v", err)
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG] Parsed start_date: %s", startDate)

	endDate, err := time.Parse("2006-01-02", payment.EndDate)
	if err != nil {
		log.Printf("[ERROR] Parsing end_date: %v", err)
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG] Parsed end_date: %s", endDate)

	// Step 1: Insert into MembershipPayment table
	log.Println("[DEBUG] Inserting membership payment into the database")
	result, err := db.Exec(`
		INSERT INTO MembershipPayment (user_id, membership_level, amount, payment_method, start_date, end_date)
		VALUES (?, ?, ?, ?, ?, ?)`,
		payment.UserID, payment.MembershipLevel, payment.Amount, payment.PaymentMethod, startDate, endDate)
	if err != nil {
		log.Printf("[ERROR] Inserting membership payment: %v", err)
		http.Error(w, "Failed to process membership payment", http.StatusInternalServerError)
		return
	}

	// Retrieve the last inserted payment ID
	paymentID, err := result.LastInsertId()
	if err != nil {
		log.Printf("[ERROR] Retrieving membership payment ID: %v", err)
		http.Error(w, "Failed to retrieve payment ID", http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] Membership payment inserted successfully. Payment ID: %d", paymentID)

	// Step 2: Call API to update the user's membership level
	log.Println("[DEBUG] Updating user membership level via API")
	apiURL := "http://localhost:5100/api/v1/user/membership/update"
	payload := map[string]interface{}{
		"user_id":         payment.UserID,
		"membership_tier": payment.MembershipLevel,
	}
	jsonPayload, _ := json.Marshal(payload)
	log.Printf("[DEBUG] Serialized JSON payload for membership update: %s", string(jsonPayload))

	req, _ := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Calling membership update API: %v", err)
		http.Error(w, "Failed to update membership level", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] Membership update API returned non-OK status: %d, Response: %s", resp.StatusCode, string(body))
		http.Error(w, "Failed to update membership level", http.StatusInternalServerError)
		return
	}
	log.Println("[DEBUG] User membership level updated successfully via API")

	// Step 3: Generate and send invoice email
	log.Println("[DEBUG] Generating and sending membership invoice email")
	err = generateMembershipInvoiceAndSendEmail(
		int(paymentID),
		payment.UserID,
		payment.MembershipLevel,
		payment.Amount,
		payment.PaymentMethod,
		payment.Email,
		startDate,
		endDate,
	)
	if err != nil {
		log.Printf("[ERROR] Sending membership invoice email: %v", err)
		http.Error(w, "Failed to send membership invoice email", http.StatusInternalServerError)
		return
	}
	log.Println("[DEBUG] Membership invoice email sent successfully")

	// Respond with a JSON object
	log.Println("[DEBUG] Sending success response to client")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Membership payment processed successfully",
		"membership_id":   paymentID,
		"membership_level": payment.MembershipLevel,
	})
	log.Println("[DEBUG] Response sent successfully")
}



func generateMembershipInvoiceAndSendEmail(paymentID int, userID int, membershipLevel string, amount float64, paymentMethod, userEmail string, startDate, endDate time.Time) error {
	// Generate the invoice in memory
	fileBytes, err := generateMembershipInvoice(paymentID, userID, membershipLevel, amount, paymentMethod, startDate, endDate)
	if err != nil {
		return fmt.Errorf("error generating invoice: %v", err)
	}

	// Email content
	subject := "Your EcoDrive Membership Invoice"
	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Invoice</title>
        </head>
        <body>
            <p>Dear User,</p>
            <p>Thank you for your EcoDrive membership purchase. Attached is your invoice:</p>
            <ul>
                <li><strong>Payment ID:</strong> %d</li>
                <li><strong>Membership Level:</strong> %s</li>
                <li><strong>Amount:</strong> $%.2f</li>
                <li><strong>Payment Method:</strong> %s</li>
                <li><strong>Start Date:</strong> %s</li>
                <li><strong>End Date:</strong> %s</li>
            </ul>
            <p>We hope you enjoy your EcoDrive membership!</p>
            <p>Best regards,<br>The EcoDrive Team</p>
        </body>
        </html>
    `, paymentID, membershipLevel, amount, paymentMethod, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Send the email with the invoice attached
	fileName := fmt.Sprintf("Membership_Invoice_%d.pdf", paymentID)
	return sendEmailWithAttachment(userEmail, subject, body, fileName, fileBytes)
}

func generateMembershipInvoice(paymentID int, userID int, membershipLevel string, amount float64, paymentMethod string, startDate, endDate time.Time) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.SetFillColor(25, 135, 84) // Green color scheme
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 10, "EcoDrive Membership Invoice", "1", 1, "C", true, 0, "")

	// Content
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Payment ID: %d", paymentID))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("User ID: %d", userID))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Membership Level: %s", membershipLevel))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Amount: $%.2f", amount))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Payment Method: %s", paymentMethod))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Start Date: %s", startDate.Format("2006-01-02")))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("End Date: %s", endDate.Format("2006-01-02")))

	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 10, "Thank you for your purchase. Enjoy your membership!", "", 1, "C", false, 0, "")

	buf := new(bytes.Buffer)
	err := pdf.Output(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}