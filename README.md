# Installation Guide

1. **Download Source Code**  
   Download the source code as a ZIP file from the provided repository link.

2. **Extract and Open**  
   Extract the ZIP file and open the project in your preferred Integrated Development Environment (IDE).

3. **Contact for Environment Files**  
   Reach out to Jeffrey via Telegram at [@cheeguang](https://t.me/cheeguang) to obtain the required environment files.

4. **Place Environment Files**  
   Place the environment files in their respective directories as instructed.

5. **Start Docker Desktop**  
   Ensure Docker Desktop is installed and running on your machine.

6. **Execute Batch File**  
   Run the `prune-build-up-loadBalance.bat` file located in the project directory to set up and initialise the application.

7. **Access the Application**  
   Open your browser and navigate to [http://127.0.0.1:8080/index.html](http://127.0.0.1:8080/index.html) to access the application.

# Overview

In the era of sustainable transportation, EcoDrive revolutionises urban mobility with a user-centric electric car-sharing platform. Designed for scalability and practicality, it features user membership tiers, promotional discounts, and precise billing. The architecture, as outlined in the diagram, ensures robust and efficient operations.

# Architecture Diagram

![Architecture Diagram](./frontend/img/CNAD-Assg1-S10258143A-ArchitectureDiagram.png)

# Functional Requirements and Architectural Rationale

## **Functional Requirements**

### **3.1. User Management**

- **3.1.1. User Registration and Authentication**  
  Implement user registration with email or phone verification. Secure authentication must include password encryption.
- **3.1.2. Membership Tiers**  
  Develop multiple membership levels (e.g., Basic, Premium, VIP) with benefits such as reduced hourly rates, priority vehicle access, or increased booking limits.
- **3.1.3. User Profile Management**  
  Enable users to update personal details, view membership status, and track rental history.

### **3.2. Vehicle Reservation System**

- **3.2.1. Car Availability and Booking**  
  Allow users to view available vehicles in real-time and make reservations for a specified time range.
- **3.2.2. Booking Modification and Cancellation**  
  Provide options to modify or cancel reservations within specified policies, with automatic updates to vehicle availability.
- **3.2.3. (BONUS) Vehicle Status Tracking**  
  Implement mechanisms to track the location, charge level, and cleanliness of vehicles to ensure readiness for the next rental.

### **3.3. Billing and Payment Processing**

- **3.3.1. Tier-Based Pricing and Discounts**  
  Calculate billing based on membership level, rental duration, and applicable promotional discounts.
- **3.3.2. Real-Time Billing Calculation**  
  Display estimated costs before booking confirmation and provide real-time cost updates during the rental period.
- **3.3.3. (BONUS) Payment Processing**  
  Integrate secure payment processing, handle refunds for cancellations, and store payment methods in compliance with industry standards.
- **3.3.4. Invoicing and Receipts**  
  Automatically generate detailed invoices after each rental and send them via email or make them accessible in the user’s profile.

## **Architectural Rationale**

### **Microservices Architecture**

EcoDrive employs a **microservices architecture** to ensure modularity, scalability, and ease of maintenance. Each microservice is responsible for a specific business functionality, enhancing flexibility and enabling independent scaling.

1. **User Microservice**:

   - **Functionality**: Manages user profiles, membership tiers, and rental history.
   - **Database**: Utilises MySQL with optimised indexing, data types, and adherence to normalisation for efficient data retrieval and minimal redundancy.

2. **Authentication Microservice**:

   - **Functionality**: Handles user registration, login processes, and secure password storage.
   - **Database**: Stores user credentials in a normalised and optimised MySQL database.
   - **Integration**: Sends email verification using Gmail SMTP for secure communication.

3. **Vehicle Microservice**:

   - **Functionality**: Tracks vehicle availability, manages bookings, modifications, and cancellations.
   - **Database**: Uses a MySQL database optimised for vehicle data storage.
   - **Integration**: Sends booking confirmations and cancellation emails through Gmail SMTP.

4. **Payment Microservice**:
   - **Functionality**: Manages billing, promotional discounts, payment processing, and invoice generation.
   - **Database**: Utilises a MySQL database optimised for transaction and financial data.
   - **Integration**:
     - **Amazon S3**: Stores invoice PDFs for administrative purposes, enabling admin access to all invoices.
     - **Gmail SMTP**: Sends refund confirmations and invoices to users.

## Technology Breakdown

<div>
  <table border="1" style="border-collapse: collapse; width: 100%; text-align: center;">
    <thead>
      <tr>
        <th><strong>Technology</strong></th>
        <th><strong>Description</strong></th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>
          <strong>Go</strong><br>
          <img src="./frontend/img/go.png" alt="Go Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td style="text-align: left;">The programming language used to develop EcoDrive's microservices, ensuring high performance, simplicity, and efficient concurrency handling.</td>
      </tr>
      <tr>
        <td>
          <strong>MySQL</strong><br>
          <img src="./frontend/img/mysql.svg" alt="MySQL Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td style="text-align: left;">
          Used for structured data storage with optimised database management practices, including:
          <ul style="text-align: left;">
            <li><strong>Indexes:</strong> Key fields are indexed to optimise query performance.</li>
            <li><strong>Normalisation:</strong> Adheres to eliminate data redundancy and maintain integrity.</li>
            <li><strong>Optimised Data Types:</strong> Improves storage efficiency and query execution.</li>
          </ul>
        </td>
      </tr>
      <tr>
        <td>
          <strong>Docker</strong><br>
          <img src="./frontend/img/docker.png" alt="Docker Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td>Used to containerise each microservice, ensuring consistency across development and production environments.</td>
      </tr>
      <tr>
        <td>
          <strong>Docker Compose</strong><br>
          <img src="./frontend/img/docker-compose.png" alt="Docker Compose Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td style="text-align: left;">
          Manages and deploys multiple interconnected containers:
          <ul style="text-align: left;">
            <li><strong>Containerisation:</strong> Each service runs in isolation.</li>
            <li><strong>Multi-Service Orchestration:</strong> Defines dependencies and networking.</li>
            <li><strong>Scalability:</strong> Allows independent scaling.</li>
            <li><strong>Environment Replication:</strong> Reduces deployment-related issues.</li>
          </ul>
        </td>
      </tr>
      <tr>
        <td>
          <strong>NGINX</strong><br>
          <img src="./frontend/img/nginx.png" alt="NGINX Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td style="text-align: left;">Acts as a load balancer to distribute traffic among frontend instances, ensuring high availability and optimal utilisation of resources.</td>
      </tr>
      <tr>
        <td>
          <strong>Gmail SMTP</strong><br>
          <img src="./frontend/img/gmail-smtp.png" alt="Gmail SMTP Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td style="text-align: left;">
          Integrated for email notifications, such as:
          <ul style="text-align: left;">
            <li><strong>Authentication Microservice:</strong> Email verification for account creation.</li>
            <li><strong>Vehicle Microservice:</strong> Booking refund confirmation.</li>
            <li><strong>Payment Microservice:</strong> Booking and membership confirmations and invoices.</li>
          </ul>
        </td>
      </tr>
      <tr>
        <td>
          <strong>Amazon S3</strong><br>
          <img src="./frontend/img/s3.png" alt="Amazon S3 Logo" style="width:50px; height:auto; margin-top: 10px;">
        </td>
        <td style="text-align: left;">Exclusively used by the Payment Microservice to store invoice PDFs, enabling administrators to access and review user transactions securely.</td>
      </tr>
    </tbody>
  </table>
</div>

# API Documentation

## Authentication Microservice

### Base URL

`http://127.0.0.1:5050/api/v1/authentication`

### **1. Send Verification Code**

#### **Endpoint**

`POST /send-verification`

#### **Description**

Sends a 6-digit verification code to the user's email and stores it in the database.

#### **Request Body**

```json
{
  "email": "user@example.com"
}
```

#### **Response**

- **200 OK**: Verification code sent successfully.
  ```json
  {
    "message": "Verification code sent successfully"
  }
  ```
- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Database or email sending error.

#### **Notes**

- If the email already exists in the database, the verification code is updated.
- Verification codes expire after 10 minutes.

### **2. Register User**

#### **Endpoint**

`POST /register-user`

#### **Description**

Registers a user after validating the provided verification code.

#### **Request Body**

```json
{
  "email": "user@example.com",
  "verification_code": "123456",
  "name": "John Doe",
  "password": "password123",
  "contact_number": "98765432",
  "address": "123 Main Street"
}
```

#### **Response**

- **200 OK**: User registered successfully.
  ```json
  {
    "message": "User registered successfully in both systems"
  }
  ```
- **400 Bad Request**: Invalid input.
- **401 Unauthorized**: Invalid or expired verification code.
- **404 Not Found**: Email not found.
- **500 Internal Server Error**: Database or microservice error.

#### **Notes**

- Password is securely hashed using `bcrypt` before storage.
- If registration is successful, the user's record is also created in the `userMicroservice`.

### **3. Authenticate User**

#### **Endpoint**

`POST /login`

#### **Description**

Authenticates a user by verifying their email and password.

#### **Request Body**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### **Response**

- **200 OK**: Login successful.
  ```json
  {
    "message": "Login successful",
    "token": "JWT_TOKEN_HERE"
  }
  ```
- **401 Unauthorized**: Invalid credentials.
- **500 Internal Server Error**: Database error.

#### **Notes**

- Password is verified using `bcrypt`.
- Returns a JWT token for authenticated sessions.

## User Microservice

### Base URL

`http://127.0.0.1:5100/api/v1/user`

## **1. Create User**

### **Endpoint**

`POST /create`

### **Description**

Creates a new user record with the default membership level set to 'Basic'.

### **Request Body**

```
{
  "name": "John Doe",
  "email": "user@example.com",
  "contact_number": "98765432",
  "address": "123 Main Street",
  "password": "hashedPassword123"
}
```

### **Response**

- **201 Created**: User created successfully.

```
{
  "message": "User created successfully"
}
```

- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Database error.

## **2. Get Membership Status**

### **Endpoint**

`GET /membership/status`

### **Description**

Retrieves the membership status of a user.

### **Query Parameters**

- `user_id` (string, required): The ID of the user.

### **Response**

- **200 OK**: Membership status retrieved successfully.

```
{
  "membership_level": "Basic"
}
```

- **404 Not Found**: User not found.
- **500 Internal Server Error**: Database error.

## **3. Update Membership Tier**

### **Endpoint**

`PUT /membership/update`

### **Description**

Updates the membership tier of a user.

### **Request Body**

```
{
  "user_id": 1,
  "membership_tier": "Premium"
}
```

### **Response**

- **200 OK**: Membership tier updated successfully.

```
{
  "message": "Membership level updated successfully"
}
```

- **400 Bad Request**: Invalid input.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Database error.

## **4. Get User Profile**

### **Endpoint**

`GET /profile`

### **Description**

Retrieves the profile information of a user.

### **Query Parameters**

- `user_id` (string, required): The ID of the user.

### **Response**

- **200 OK**: User profile retrieved successfully.

```
{
  "name": "John Doe",
  "email": "user@example.com",
  "contact_number": "98765432",
  "address": "123 Main Street",
  "membership_level": "Basic"
}
```

- **404 Not Found**: User not found.
- **500 Internal Server Error**: Database error.

## **5. Update User Profile**

### **Endpoint**

`PUT /profile/update`

### **Description**

Updates the profile information of a user.

### **Request Body**

```
{
  "user_id": 1,
  "name": "John Doe",
  "contact_number": "98765432",
  "address": "123 Main Street"
}
```

### **Response**

- **200 OK**: User profile updated successfully.

```
{
  "message": "Profile updated successfully"
}
```

- **400 Bad Request**: Invalid input.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Database error.

## **6. Get Rental History**

### **Endpoint**

`GET /rental-history`

### **Description**

Retrieves the rental history of a user.

### **Query Parameters**

- `user_id` (string, required): The ID of the user.

### **Response**

- **200 OK**: Rental history retrieved successfully.

```
[
  {
    "vehicle_id": 101,
    "rental_price_per_hour": 15.0,
    "rental_date": "2024-01-15T10:30:00Z"
  },
  {
    "vehicle_id": 102,
    "rental_price_per_hour": 20.0,
    "rental_date": "2024-02-10T12:45:00Z"
  }
]
```

- **404 Not Found**: User not found.
- **500 Internal Server Error**: Database error.

## Vehicle Microservice

### Base URL

`http://127.0.0.1:5150/api/v1/vehicle`

## **1. Get Available Vehicles**

### **Endpoint**

`GET /availability`

### **Description**

Retrieves a list of available vehicles for the specified date range.

### **Query Parameters**

- `start_date` (string, required): The start date in ISO 8601 format (`YYYY-MM-DDTHH:MM`).
- `end_date` (string, required): The end date in ISO 8601 format (`YYYY-MM-DDTHH:MM`).

### **Response**

- **200 OK**: List of available vehicles retrieved successfully.

```
[
  {
    "vehicle_id": 1,
    "model": "Tesla Model S",
    "location": "Downtown",
    "charge_level": 80,
    "cleanliness_status": "Clean",
    "rental_price_per_hour": 50.0
  }
]
```

- **400 Bad Request**: Missing or invalid query parameters.
- **500 Internal Server Error**: Database error.

## **2. Get Vehicle Status**

### **Endpoint**

`GET /status`

### **Description**

Retrieves the status of all vehicles, including their availability, location, and condition.

### **Response**

- **200 OK**: List of vehicles and their statuses retrieved successfully.

```
[
  {
    "vehicle_id": 1,
    "model": "Tesla Model S",
    "location": "Downtown",
    "charge_level": 80,
    "cleanliness_status": "Clean",
    "rental_price_per_hour": 50.0
  }
]
```

- **500 Internal Server Error**: Database error.

## **3. Update Vehicle Details**

### **Endpoint**

`PUT /update`

### **Description**

Updates the details of a specific vehicle, including cleanliness status, charge level, and location. Optionally updates the booking status.

### **Request Body**

```
{
  "vehicle_id": 1,
  "cleanliness_status": "Clean",
  "charge_level": 90,
  "location": "Downtown Parking Lot",
  "booking_id": 12345
}
```

### **Response**

- **200 OK**: Vehicle details updated successfully.

```
{
  "message": "Vehicle details and booking status updated successfully"
}
```

- **400 Bad Request**: Missing or invalid input fields.
- **404 Not Found**: Vehicle not found or no changes made.
- **500 Internal Server Error**: Database error.

## **4. Create Booking**

### **Endpoint**

`POST /booking`

### **Description**

Creates a new booking for a vehicle.

### **Request Body**

```
{
  "vehicle_id": 1,
  "user_id": 123,
  "booking_date": "2024-12-01T10:00",
  "return_date": "2024-12-03T18:00",
  "total_price": 200.0
}
```

### **Response**

- **201 Created**: Booking created successfully.

```
{
  "booking_id": 1001
}
```

- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Database error.

## **5. Get Booking**

### **Endpoint**

`GET /booking/{id}`

### **Description**

Retrieves the details of a specific booking by ID.

### **Path Parameters**

- `id` (int, required): The ID of the booking.

### **Response**

- **200 OK**: Booking details retrieved successfully.

```
{
  "booking_id": 1001,
  "vehicle_id": 1,
  "user_id": 123,
  "booking_date": "2024-12-01T10:00",
  "return_date": "2024-12-03T18:00",
  "total_price": 200.0,
  "model": "Tesla Model S",
  "location": "Downtown",
  "charge_level": 80
}
```

- **400 Bad Request**: Invalid booking ID.
- **404 Not Found**: Booking not found.
- **500 Internal Server Error**: Database error.

## **6. Modify Booking**

### **Endpoint**

`PUT /booking/{id}`

### **Description**

Modifies the details of an existing booking.

### **Path Parameters**

- `id` (int, required): The ID of the booking.

### **Request Body**

```
{
  "start_date_time": "2024-12-01T10:00",
  "end_date_time": "2024-12-03T18:00",
  "total_price": 250.0
}
```

### **Response**

- **200 OK**: Booking updated successfully.

```
{
  "message": "Booking updated successfully"
}
```

- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Database error.

## **7. Cancel Booking**

### **Endpoint**

`DELETE /booking/{id}`

### **Description**

Cancels an existing booking and sends an acknowledgment email to the user.

### **Path Parameters**

- `id` (int, required): The ID of the booking.

### **Request Body**

```
{
  "email": "user@example.com"
}
```

### **Response**

- **200 OK**: Booking cancelled successfully, and acknowledgment email sent.

```
{
  "message": "Booking cancelled successfully"
}
```

- **400 Bad Request**: Invalid input.
- **500 Internal Server Error**: Database or email sending error.

## **8. Get Bookings by User**

### **Endpoint**

`GET /booking/user/{user_id}`

### **Description**

Retrieves all bookings for a specific user.

### **Path Parameters**

- `user_id` (int, required): The ID of the user.

### **Response**

- **200 OK**: List of bookings retrieved successfully.

```
[
  {
    "booking_id": 1001,
    "vehicle_id": 1,
    "user_id": 123,
    "booking_date": "2024-12-01T10:00",
    "return_date": "2024-12-03T18:00",
    "total_price": 200.0,
    "model": "Tesla Model S",
    "location": "Downtown",
    "charge_level": 80,
    "rental_price_per_hour": 50.0,
    "cleanliness_status": "Clean",
    "status": "active"
  }
]
```

- **400 Bad Request**: Invalid user ID.
- **500 Internal Server Error**: Database error.

## **9. Get Bookings by Vehicle**

### **Endpoint**

`GET /booking/vehicle/{vehicle_id}`

### **Description**

Retrieves all booking dates for a specific vehicle.

### **Path Parameters**

- `vehicle_id` (int, required): The ID of the vehicle.

### **Response**

- **200 OK**: List of booking dates retrieved successfully.

```
[
  {
    "booking_date": "2024-12-01T10:00",
    "return_date": "2024-12-03T18:00"
  }
]
```

- **400 Bad Request**: Invalid vehicle ID.
- **500 Internal Server Error**: Database error.

## **10. End Booking**

### **Endpoint**

`PUT /booking/end/{id}`

### **Description**

Marks a booking as completed and updates the associated vehicle details.

### **Path Parameters**

- `id` (int, required): The ID of the booking.

### **Request Body**

```
{
  "charge_level": 100,
  "cleanliness_status": "Clean",
  "location": "Downtown Parking Lot"
}
```

### **Response**

- **200 OK**: Booking ended successfully.

```
{
  "message": "Booking ended successfully"
}
```

- **400 Bad Request**: Missing or invalid input.
- **500 Internal Server Error**: Database error.

## Payment Microservice

### Base URL

`http://127.0.0.1:5200/api/v1/payment`

## **1. Calculate Real-Time Bill**

### **Endpoint**

`GET /real-time-bill`

### **Description**

Calculates the total cost, discount, and final price based on membership level and rental duration.

### **Query Parameters**

- `membership_level` (string, required): The user's membership level.
- `duration_hours` (integer, required): Duration of the rental in hours.
- `price_per_hour` (float, required): Rental cost per hour.

### **Response**

- **200 OK**: Real-time bill calculated successfully.

```
{
  "final_price": 180.0,
  "discount": 20.0,
  "total_price": 200.0,
  "membership": "Gold",
  "duration": 10
}
```

- **400 Bad Request**: Invalid input or missing parameters.
- **500 Internal Server Error**: Database error.

## **2. Process Payment**

### **Endpoint**

`POST /process`

### **Description**

Processes a payment for a vehicle rental, notifies the booking service, and generates an invoice.

### **Request Body**

```
{
  "user_id": 123,
  "vehicle_id": "101",
  "start_date": "2024-12-01T10:00",
  "end_date": "2024-12-03T18:00",
  "payment_method": "Credit Card",
  "price_per_hour": "20.0",
  "rental_duration": "2",
  "total_price": "200.0",
  "email": "user@example.com"
}
```

### **Response**

- **200 OK**: Payment processed successfully, booking service notified, and invoice sent.

```
{
  "message": "Payment processed successfully",
  "booking_id": 1001,
  "payment_id": 2001
}
```

- **400 Bad Request**: Invalid input or missing parameters.
- **500 Internal Server Error**: Database error, booking service error, or email failure.

## **3. Process Membership Payment**

### **Endpoint**

`POST /membership/payment`

### **Description**

Processes a payment for a membership plan, updates the user's membership level, and generates an invoice.

### **Request Body**

```
{
  "user_id": 123,
  "membership_level": "Gold",
  "amount": 100.0,
  "payment_method": "PayNow",
  "start_date": "2024-12-01",
  "end_date": "2025-12-01",
  "email": "user@example.com"
}
```

### **Response**

- **200 OK**: Membership payment processed, user membership updated, and invoice sent.

```
{
  "message": "Membership payment processed successfully",
  "membership_id": 3001,
  "membership_level": "Gold"
}
```

- **400 Bad Request**: Invalid input or missing parameters.
- **401 Unauthorized**: Missing or invalid authorization token.
- **500 Internal Server Error**: Database error or membership service error.
