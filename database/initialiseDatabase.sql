-- **************************************************
-- DATABASE: ecoDrive_authentication_db
-- PURPOSE: Handles authentication and registration 
--          information for the Authentication Service
-- **************************************************

-- Drop and recreate the authentication database
DROP DATABASE IF EXISTS ecoDrive_authentication_db;
CREATE DATABASE ecoDrive_authentication_db;
USE ecoDrive_authentication_db;

-- Create the Authentication table
-- PURPOSE: Stores authentication tokens and expiry information
CREATE TABLE Authentication (
    auth_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for authentication record
    user_id SMALLINT UNSIGNED NOT NULL,                             -- Associated user ID
    auth_token VARCHAR(500) NOT NULL UNIQUE,                        -- Authentication token
    token_expiry TIMESTAMP NOT NULL,                                -- Expiry timestamp of the token
    INDEX idx_user_id (user_id),                                    -- Index to optimise lookups by user_id
    INDEX idx_token_expiry (token_expiry)                           -- Index to optimise expiry checks
);

-- Insert example data into the Authentication table
INSERT INTO Authentication (user_id, auth_token, token_expiry) VALUES
(1, "random_generated_token", "2024-12-31 23:59:59");


-- Create the User table
-- PURPOSE: Stores user registration details (previously named Registration)
CREATE TABLE User (
    user_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for the user
    name VARCHAR(100),                                              -- User's name
    email VARCHAR(100) NOT NULL UNIQUE,                            -- User's email address
    password VARCHAR(255) NULL,                                    -- User's password (hashed)
    contact_number VARCHAR(15),                                    -- User's contact number
    address TEXT,                                                  -- User's address
    verification_code VARCHAR(10) NOT NULL,                        -- Verification code
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                -- Record creation timestamp
    INDEX idx_email (email),                                       -- Index to optimise lookups by email
    INDEX idx_created_at (created_at)                              -- Index to optimise recent user lookups
);

-- Insert example data into the User table
INSERT INTO User (name, email, password, contact_number, address, verification_code, created_at) VALUES
("John Tan", "john@gmail.com", "$2a$10$LEL8btFg0WcO7BaUBI5JJ.3kqEv/hv6s1bM2u6DV0to71cIkDaadK", "12345678", "123 Main St, Singapore", "123456", NOW());


-- **************************************************
-- DATABASE: ecoDrive_user_db
-- PURPOSE: Stores user profile information for the 
--          User Service
-- **************************************************

-- Drop and recreate the user database
DROP DATABASE IF EXISTS ecoDrive_user_db;
CREATE DATABASE ecoDrive_user_db;
USE ecoDrive_user_db;

-- Create the User table
-- PURPOSE: Stores user personal and profile details
CREATE TABLE User (
    user_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for the user
    name VARCHAR(100),                                              -- User's name
    email VARCHAR(100) NOT NULL UNIQUE,                            -- User's email address
    password VARCHAR(255) NOT NULL,                                -- User's password (hashed)
    membership_level ENUM('Basic', 'Premium', 'VIP'),              -- Membership tier
    contact_number VARCHAR(15),                                    -- User's contact number
    address TEXT,                                                  -- User's address
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                -- Record creation timestamp
    INDEX idx_email_membership (email, membership_level)           -- Composite index for email and membership level
);

-- Insert example data into the User table
INSERT INTO User (name, email, password, membership_level, contact_number, address) VALUES
("John Tan", "john@gmail.com", "$2a$10$LEL8btFg0WcO7BaUBI5JJ.3kqEv/hv6s1bM2u6DV0to71cIkDaadK", "Basic", "12345678", "123 Main St, Singapore");


-- **************************************************
-- DATABASE: ecoDrive_vehicle_db
-- PURPOSE: Tracks vehicle details and availability 
--          for the Vehicle Service
-- **************************************************

-- Drop and recreate the vehicle database
DROP DATABASE IF EXISTS ecoDrive_vehicle_db;
CREATE DATABASE ecoDrive_vehicle_db;
USE ecoDrive_vehicle_db;

-- Create the Vehicles table
-- PURPOSE: Stores vehicle details and rental information
CREATE TABLE Vehicles (
    vehicle_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT, -- Unique ID for the vehicle
    model VARCHAR(255) NOT NULL,                                     -- Model of the vehicle
    location VARCHAR(255),                                           -- Current location of the vehicle
    charge_level TINYINT UNSIGNED,                                   -- Battery charge level (for EVs)
    cleanliness_status ENUM('Clean', 'Needs Cleaning'),             -- Cleanliness status of the vehicle
    rental_price_per_hour DECIMAL(10, 2) NOT NULL,                   -- Rental price per hour
    INDEX idx_location (location),                                  -- Index for location-based searches
    INDEX idx_charge_level (charge_level)                           -- Index for charge level lookups
);

-- Create the Bookings table
-- PURPOSE: Stores booking details linked to vehicles
CREATE TABLE Bookings (
    booking_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT, -- Unique ID for the booking
    vehicle_id SMALLINT UNSIGNED NOT NULL,                            -- Vehicle ID (foreign key)
    user_id SMALLINT UNSIGNED NOT NULL,                               -- User ID (not a foreign key)
    booking_date DATETIME NOT NULL,                                   -- Date and time of booking
    return_date DATETIME NOT NULL,                                    -- Date and time of return
    total_price DECIMAL(10, 2) NOT NULL,                              -- Total price of the booking
    FOREIGN KEY (vehicle_id) REFERENCES Vehicles(vehicle_id),         -- Foreign key relationship
    INDEX idx_user_booking_date (user_id, booking_date)               -- Composite index for user and booking date
);

-- Insert example data into the Vehicles table
INSERT INTO Vehicles (model, location, charge_level, cleanliness_status, rental_price_per_hour) VALUES
("Toyota Prius", "Marina Barrage Public Carpark", 95, "Clean", 25.00),
("Tesla Model 3", "ION Orchard Car Park", 80, "Needs Cleaning", 50.00),
("Honda Civic", "NEX Carpark", 0, "Clean", 20.00),
("Nissan Leaf", "Prime Auto Care VivoCity (Yellow Zone) B2 Carpark", 100, "Clean", 30.00),
("Ford Mustang", "Suntec City Carpark F", 0, "Needs Cleaning", 40.00);

-- Insert example data into the Bookings table
INSERT INTO Bookings (vehicle_id, user_id, booking_date, return_date, total_price) VALUES
(1, 1, '2025-01-01 10:00:00', '2025-01-05 14:00:00', 100.00),
(1, 1, '2024-06-01 10:00:00', '2024-06-05 14:00:00', 100.00),
(1, 1, '2024-01-01 10:00:00', '2024-01-05 14:00:00', 100.00);


-- **************************************************
-- DATABASE: ecoDrive_payment_db
-- PURPOSE: Tracks payment transactions for the 
--          Payment Service
-- **************************************************

-- Drop and recreate the payment database
DROP DATABASE IF EXISTS ecoDrive_payment_db;
CREATE DATABASE ecoDrive_payment_db;
USE ecoDrive_payment_db;

-- Rename Payments to BookingPayment
-- PURPOSE: Tracks payments related to bookings
CREATE TABLE BookingPayment (
    payment_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for payment
    user_id SMALLINT UNSIGNED NOT NULL,                                -- Associated user ID
    booking_id SMALLINT UNSIGNED NOT NULL,                             -- Booking reference ID
    amount DECIMAL(10, 2) NOT NULL,                                    -- Payment amount
    payment_method ENUM('Card', 'PayNow'),                             -- Payment method used
    payment_status ENUM('Pending', 'Completed', 'Refunded'),           -- Status of the payment
    discount DECIMAL(10, 2) DEFAULT 0.00,                              -- Discount amount
    final_amount DECIMAL(10, 2) DEFAULT 0.00,                          -- Final amount after discount
    invoice_pdf TEXT,                                                  -- Path to the invoice PDF
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                    -- Record creation timestamp
    INDEX idx_user_payment_status (user_id, payment_status)            -- Composite index for user and payment status
);

-- Insert example data into the BookingPayment table
INSERT INTO BookingPayment (user_id, booking_id, amount, payment_method, payment_status, discount, final_amount) VALUES
(1, 101, 100.50, "Card", "Completed", 10.00, 90.50);

-- Create the MembershipPayment table
-- PURPOSE: Tracks payments related to membership plans
CREATE TABLE MembershipPayment (
    membership_payment_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT, -- Unique ID for membership payment
    user_id SMALLINT UNSIGNED NOT NULL,                                -- Associated user ID
    membership_level ENUM('Basic', 'Premium', 'VIP') NOT NULL,         -- Membership tier purchased
    amount DECIMAL(10, 2) NOT NULL,                                    -- Payment amount for membership
    payment_method ENUM('Card', 'PayNow'),                             -- Payment method used
    payment_status ENUM('Pending', 'Completed', 'Refunded'),           -- Status of the payment
    start_date DATE NOT NULL,                                          -- Membership start date
    end_date DATE NOT NULL,                                            -- Membership end date
    invoice_pdf TEXT,                                                  -- Path to the invoice PDF
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                    -- Record creation timestamp
    INDEX idx_user_membership_level (user_id, membership_level)        -- Composite index for user and membership level
);

-- Insert example data into the MembershipPayment table
INSERT INTO MembershipPayment (user_id, membership_level, amount, payment_method, payment_status, start_date, end_date) VALUES
(2, 'Premium', 150.00, "Card", "Completed", '2024-01-01', '2024-12-31');

-- Create the Discounts table
-- PURPOSE: Manages promotional discounts by membership level
CREATE TABLE Discounts (
    discount_id SMALLINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT, -- Unique ID for the discount
    membership_level ENUM('Basic', 'Premium', 'VIP') NOT NULL,         -- Membership tier
    discount_percentage DECIMAL(5, 2) NOT NULL,                       -- Discount percentage
    INDEX idx_membership_discount (membership_level)                  -- Index for quick lookup by membership level
);

-- Insert example data into the Discounts table
INSERT INTO Discounts (membership_level, discount_percentage) VALUES
('Basic', 5.00),
('Premium', 10.00),
('VIP', 20.00);
