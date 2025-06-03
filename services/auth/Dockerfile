# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.mod file to speed up builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# Specifically build the auth service main.go. Adjust path if needed.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./auth-service ./services/auth/main.go

# Start a new stage from scratch for a smaller image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/auth-service .
COPY .env .env

# Expose port (though gRPC doesn't strictly need EXPOSE in Dockerfile, it's good for documentation)
# The actual port will be mapped in docker-compose.yml
# EXPOSE 8004

# Command to run the executable
CMD ["./auth-service"]
