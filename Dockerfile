# Use the official Golang image to build the Go application
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and the application code to the container
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the application binary
RUN go build -o main ./cmd/server/main.go

# Expose port 8080 for the API
ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
