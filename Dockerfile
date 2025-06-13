# Use official Go 1.18 Alpine image as base (adjust version/language if different)
FROM golang:1.18-alpine

# Set working directory
WORKDIR /app

# Copy the application code
COPY . .

# Build the Go application
RUN go build -o main

# Expose port 8080 (matches WEBSITES_PORT in Terraform)
EXPOSE 8080

# Run the application
CMD ["./main"]