package registration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Registration represents the structure of the data in the registration table.
type Registration struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	VerificationCode string `json:"verification_code"`
	CreatedAt        string `json:"created_at"`
}

// DB connection details
var db *sql.DB

func init() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	log.Println("Initializing database connection...")
	db, err = sql.Open("mysql", os.Getenv("DB_CONNECTION"))
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

// SendVerificationCode handles sending the verification code and storing it in the database.
func SendVerificationCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /send-verification request...")
	var registration Registration

	// Parse the incoming request
	err := json.NewDecoder(r.Body).Decode(&registration)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed request: %+v", registration)

	// Generate a random 6-digit verification code
	rand.Seed(time.Now().UnixNano())
	verificationCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	log.Printf("Generated verification code: %s", verificationCode)

	// Insert or update email and verification code in the database
	log.Println("Inserting or updating verification code in the database...")
	_, err = db.Exec(`
		INSERT INTO registration (email, verification_code, created_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		verification_code = VALUES(verification_code), created_at = VALUES(created_at)
	`, registration.Email, verificationCode, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("Error inserting or updating database: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	log.Println("Verification code inserted or updated in the database.")

	// Send the verification code via email
	log.Printf("Sending verification code to %s...", registration.Email)
	err = sendEmail(registration.Email, verificationCode)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}
	log.Println("Verification code sent successfully.")

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Verification code sent successfully"}`))
}


// sendEmail sends an email containing the verification code.
func sendEmail(to, code string) error {
	// SMTP configuration from .env
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Email content (HTML)
	from := "EcoDrive <" + smtpUser + ">"
	subject := "Your Verification Code"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Verification Code</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f9f9f9;
					color: #333;
					margin: 0;
					padding: 0;
				}
				.container {
					width: 100%%;
					max-width: 600px;
					margin: 0 auto;
					background: #ffffff;
					padding: 20px;
					border-radius: 10px;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #4CAF50;
				}
				.code {
					font-size: 20px;
					font-weight: bold;
					color: #4CAF50;
					margin: 20px 0;
				}
				.footer {
					margin-top: 20px;
					font-size: 12px;
					color: #777;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>EcoDrive Verification Code</h1>
				<p>Dear User,</p>
				<p>Thank you for signing up with EcoDrive! Please use the following verification code to complete your registration:</p>
				<div class="code">%s</div>
				<p>If you did not request this email, please ignore it.</p>
				<p>Best regards,</p>
				<p>The EcoDrive Team</p>
				<div class="footer">
					<p>EcoDrive &copy; 2024. All Rights Reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, code)

	// Combine headers and body
	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body)

	// Authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	log.Println("Email sent successfully.")
	return nil
}