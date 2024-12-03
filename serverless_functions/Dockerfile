# Use an official Golang image as a builder
FROM golang:latest as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

RUN ls -la /app
RUN ls -la /app/cmd

# Build the application
RUN go build -o request-handler ./cmd/controller

# Use a minimal image for the runtime
FROM gcr.io/distroless/static:nonroot

# Set the working directory inside the container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/request-handler .

# Expose the port that the function listens on
EXPOSE 8080

# Set the entry point
CMD ["/request-handler"]