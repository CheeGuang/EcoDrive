# Use an official lightweight Go image
FROM golang:1.23-alpine

# Install bash as it is required for the wait-for-it.sh script
RUN apk add --no-cache bash

# Set environment variables for Go
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Copy the wait-for-it script
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Build the Go application
RUN go build -o vehicleMicroservice .

# Expose the service port
EXPOSE 5150

# Command to wait for MySQL and run the application
CMD ["/bin/bash", "/wait-for-it.sh", "mysql:3306", "-t", "60", "--", "./vehicleMicroservice"]
