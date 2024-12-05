# Use an official Golang image as a builder
FROM golang:latest as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pod_creation_handler ./cmd/controller

# Use a minimal image for the runtime
FROM gcr.io/distroless/static:nonroot

# Set the working directory inside the container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/pod_creation_handler /pod_creation_handler

# Expose the port that the function listens on
EXPOSE 8080

# Set the entry point
CMD ["/pod_creation_handler"]