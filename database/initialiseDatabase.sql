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
    auth_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for authentication record
    user_id INT NOT NULL,                             -- Associated user ID
    auth_token VARCHAR(1000) NOT NULL UNIQUE,          -- Authentication token
    token_expiry TIMESTAMP NOT NULL                   -- Expiry timestamp of the token
);

-- Insert example data into the Authentication table
INSERT INTO Authentication (user_id, auth_token, token_expiry) VALUES
(1, "random_generated_token", "2024-12-31 23:59:59");


-- Create the User table
-- PURPOSE: Stores user registration details (previously named Registration)
CREATE TABLE User (
    user_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for the user
    name VARCHAR(100),                                -- User's name
    email VARCHAR(100) NOT NULL UNIQUE,              -- User's email address
    password VARCHAR(255) NULL,                      -- User's password (hashed)
    contact_number VARCHAR(15),                      -- User's contact number
    address TEXT,                                    -- User's address
    verification_code VARCHAR(10) NOT NULL,          -- Verification code
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Record creation timestamp
);

-- Insert example data into the User table
INSERT INTO User (name, email, password, membership_level, contact_number, address, verification_code, created_at) VALUES
("Alice Johnson", "alice.johnson@example.com", "hashed_password", "12345678", "123 Main St, Singapore", "123456", NOW());


-- **************************************************
-- DATABASE: ecoDrive_payment_db
-- PURPOSE: Tracks payment transactions for the 
--          Payment Service
-- **************************************************

-- Drop and recreate the payment database
DROP DATABASE IF EXISTS ecoDrive_payment_db;
CREATE DATABASE ecoDrive_payment_db;
USE ecoDrive_payment_db;

-- Create the Payments table
-- PURPOSE: Records payment details and statuses
CREATE TABLE Payments (
    payment_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for payment
    user_id INT NOT NULL,                                -- Associated user ID
    booking_id INT NOT NULL,                             -- Booking reference ID
    amount DECIMAL(10, 2) NOT NULL,                      -- Payment amount
    payment_method ENUM('Card', 'PayNow'),               -- Payment method used
    payment_status ENUM('Pending', 'Completed', 'Refunded'), -- Status of the payment
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP       -- Record creation timestamp
);

-- Insert example data into the Payments table
INSERT INTO Payments (user_id, booking_id, amount, payment_method, payment_status) VALUES
(1, 101, 100.50, "Card", "Completed");


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
    user_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for the user
    name VARCHAR(100),                                -- User's name
    email VARCHAR(100) NOT NULL UNIQUE,              -- User's email address
    password VARCHAR(255) NOT NULL,                  -- User's password (hashed)
    membership_level ENUM('Basic', 'Premium', 'VIP'),-- Membership tier
    contact_number VARCHAR(15),                      -- User's contact number
    address TEXT,                                    -- User's address
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Record creation timestamp
);

-- Insert example data into the User table
INSERT INTO User (name, email, password, membership_level, contact_number, address) VALUES
("Alice Johnson", "alice.johnson@example.com", "hashed_password", "Premium", "12345678", "123 Main St, Singapore");


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
    vehicle_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,  -- Unique ID for the vehicle
    vehicle_type ENUM('Car', 'Bike', 'EV') NOT NULL,     -- Type of vehicle
    availability_status ENUM('Available', 'Booked', 'Maintenance') NOT NULL, -- Availability status
    location VARCHAR(255),                               -- Current location of the vehicle
    charge_level INT,                                    -- Battery charge level (for EVs)
    cleanliness_status ENUM('Clean', 'Needs Cleaning'), -- Cleanliness status of the vehicle
    rental_price_per_hour DECIMAL(10, 2) NOT NULL        -- Rental price per hour
);

-- Insert example data into the Vehicles table
INSERT INTO Vehicles (vehicle_type, availability_status, location, charge_level, cleanliness_status, rental_price_per_hour) VALUES
("Car", "Available", "Downtown Lot A", 100, "Clean", 20.00);
