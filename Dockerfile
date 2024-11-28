# Use the official Golang image to build the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /src

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o /app/main main.go

# Use a minimal image for the final container
FROM alpine:latest

# Install dumb-init
RUN apk --no-cache add dumb-init

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 8080

# Use dumb-init as the entrypoint to handle PID 1 correctly
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Command to run the Go application
CMD ["/app/main"]