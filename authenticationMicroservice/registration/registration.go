package registration

import (
	"bytes"
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
	"golang.org/x/crypto/bcrypt"
)

// User represents the structure of the data in the User table.
type User struct {
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
	var user User

	// Parse the incoming request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed request: %+v", user)

	// Generate a random 6-digit verification code
	rand.Seed(time.Now().UnixNano())
	verificationCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	log.Printf("Generated verification code: %s", verificationCode)

	// Insert or update email and verification code in the database
	log.Println("Inserting or updating verification code in the database...")
	_, err = db.Exec(`
		INSERT INTO User (email, verification_code, created_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		verification_code = VALUES(verification_code), created_at = VALUES(created_at)
	`, user.Email, verificationCode, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("Error inserting or updating database: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	log.Println("Verification code inserted or updated in the database.")

	// Send the verification code via email
	log.Printf("Sending verification code to %s...", user.Email)
	err = sendEmail(user.Email, verificationCode)
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

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /register-user request...")

	// Parse the incoming request
	var user struct {
		Email           string `json:"email"`
		VerificationCode string `json:"verification_code"`
		Name            string `json:"name"`
		Password        string `json:"password"`
		ContactNumber   string `json:"contact_number"`
		Address         string `json:"address"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed request: %+v", user)

	// Verify the provided verification code and timestamp
	var dbVerificationCode string
	var dbCreatedAt string

	log.Println("Checking verification code in the database...")
	err = db.QueryRow(
		"SELECT verification_code, created_at FROM User WHERE email = ?",
		user.Email,
	).Scan(&dbVerificationCode, &dbCreatedAt)
	if err == sql.ErrNoRows {
		log.Println("Email not found.")
		http.Error(w, "Email not found or verification code invalid", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	log.Printf("Database values: verification_code=%s, created_at=%s", dbVerificationCode, dbCreatedAt)

	// Convert dbCreatedAt from string to time.Time
	createdAt, err := time.Parse("2006-01-02 15:04:05", dbCreatedAt)
	if err != nil {
		log.Printf("Error parsing created_at timestamp: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if the verification code matches
	if dbVerificationCode != user.VerificationCode {
		log.Println("Verification code mismatch.")
		http.Error(w, "Invalid verification code", http.StatusUnauthorized)
		return
	}

	// Check if the code is still valid (within 10 minutes)
	if time.Since(createdAt) > 10*time.Minute {
		log.Println("Verification code expired.")
		http.Error(w, "Verification code expired", http.StatusUnauthorized)
		return
	}

	// Hash the password
	log.Println("Hashing the password...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Password hashed successfully.")

	// Update the local User table
	log.Println("Updating local User table...")
	_, err = db.Exec(`
		UPDATE User
		SET name = ?, password = ?, contact_number = ?, address = ?
		WHERE email = ?`,
		user.Name, string(hashedPassword), user.ContactNumber, user.Address, user.Email)
	if err != nil {
		log.Printf("Error updating local User table: %v", err)
		http.Error(w, "Failed to register user locally", http.StatusInternalServerError)
		return
	}
	log.Println("Local User table updated successfully.")

	// Call the userMicroservice to add the user record
	log.Println("Calling userMicroservice to create the user record...")
	err = callUserMicroservice(user.Name, user.Email, user.ContactNumber, user.Address, string(hashedPassword))
	if err != nil {
		log.Printf("Error calling userMicroservice: %v", err)
		http.Error(w, "Failed to register user in userMicroservice", http.StatusInternalServerError)
		return
	}

	log.Println("User registration successful.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User registered successfully in both systems"}`))
}


// callUserMicroservice sends a request to userMicroservice to create a new user record
func callUserMicroservice(name, email, contactNumber, address, password string) error {
	userMicroserviceURL := "http://localhost:5100/api/v1/user/create" // Update with the actual userMicroservice URL

	// Create the request payload
	payload := map[string]string{
		"name":           name,
		"email":          email,
		"contact_number": contactNumber,
		"address":        address,
		"password":       password,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return err
	}

	// Create the HTTP POST request
	resp, err := http.Post(userMicroserviceURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error sending request to userMicroservice: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		log.Printf("Unexpected status from userMicroservice: %v", resp.Status)
		return fmt.Errorf("failed to create user record, status: %v", resp.Status)
	}

	log.Println("User record created successfully in userMicroservice.")
	return nil
}
