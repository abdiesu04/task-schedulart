# Stage 1: Build Stage
# Use golang:1.21-alpine as base image - Alpine is lightweight and perfect for Go apps
# This is a multi-stage build to keep the final image size small
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
# All subsequent commands will run from this directory
WORKDIR /app

# Copy only the go.mod file first
# This is done separately to leverage Docker's layer caching
# If go.mod hasn't changed, Docker will use cached dependencies
COPY go.mod ./

# Download all dependencies specified in go.mod
# This layer will be cached if go.mod doesn't change
RUN go mod download

# Copy the entire source code into the container
# This includes all .go files and other assets
COPY . .

# Build the Go application
# CGO_ENABLED=0 - Disables C Go, creating a static binary
# GOOS=linux - Builds for Linux OS
# -o main - Names the output binary 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Final Stage
# Use alpine:latest as the base image for the final container
# This creates an extremely small final image
FROM alpine:latest

# Set working directory in the final container
WORKDIR /app

# Copy only the binary from the builder stage
# This is the magic of multi-stage builds - we only take what we need
COPY --from=builder /app/main .

# Document that the container listens on port 8080
# This is for documentation - it doesn't actually open the port
EXPOSE 8080

# Command to run when the container starts
# Using array syntax which is the preferred method
CMD ["./main"] 