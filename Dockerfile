# Use the Go 1.23 image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application
COPY . .

# Build the Go application
RUN go build -o /auth-service

# Expose the port and run the application
EXPOSE 8080
CMD ["/auth-service"]
