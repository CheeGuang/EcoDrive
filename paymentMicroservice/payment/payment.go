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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
			INSERT INTO Payments (user_id, booking_id, amount, payment_method, payment_status, discount, final_amount)
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


// uploadInvoiceToS3 uploads the invoice PDF to S3 and returns the object URL
func uploadInvoiceToS3(filePath string, fileName string) (string, error) {
	// Initialize AWS SDK
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"), // Replace with your region
	})
	if err != nil {
		log.Printf("Error initializing AWS session: %v", err)
		return "", err
	}

	s3Client := s3.New(sess)
	bucketName := "cnad-ecodrive"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file for S3 upload: %v", err)
		return "", err
	}
	defer file.Close()

	// Upload the file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}

	// Generate a presigned URL for the uploaded file
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	presignedURL, err := req.Presign(15 * time.Minute) // URL valid for 15 minutes
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
		return "", err
	}

	log.Printf("Invoice uploaded to S3. Presigned URL: %s", presignedURL)
	return presignedURL, nil
}

// sendEmail sends an email with the given content
func sendEmail(to, subject, body string) error {
	// SMTP configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	from := "EcoDrive <" + smtpUser + ">"

	// Create email
	message := bytes.NewBuffer(nil)
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)

	// Send email
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, message.Bytes())
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Println("Email sent successfully.")
	return nil
}
// encodeToBase64 encodes bytes to a base64 string
func encodeToBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}
