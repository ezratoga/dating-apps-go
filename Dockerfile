# Stage 1: Build the application
FROM golang:1.22.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file if your app requires it (uncomment if needed)
# COPY .env .env

# Expose the application port (change if the app uses a different port)
EXPOSE 8888

# Run the application
CMD ["./main"]
