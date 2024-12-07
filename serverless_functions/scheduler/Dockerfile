# Step 1: Build stage
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Step 2: Copy the go mod and sum files
COPY go.mod go.sum ./

# Step 3: Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Step 4: Copy the source code into the container
COPY . .

# Step 5: Build the Go app
RUN go build -o /scheduler cmd/scheduler/main.go

# Step 6: Start a new stage from distroless base image for the final image
FROM gcr.io/distroless/base

# Step 7: Copy the binary from the builder stage
COPY --from=builder /scheduler /scheduler

# Step 8: Expose port 8081 (where your app listens)
EXPOSE 8080

# Step 9: Set the environment variable for the port if needed
#ENV PORT=8080
#ENV DOMAIN=localhost
#ENV ETCD_ENDPOINTS=localhost:2379

# Step 10: Run the Go binary
CMD ["/scheduler"]