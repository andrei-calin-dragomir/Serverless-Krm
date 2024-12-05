# Use an official Golang image as a builder
FROM golang:latest as builder

ARG TARGETOS
ARG TARGETARCH

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . ./

# Install dependencies and tidy up the go.mod and go.sum files.
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -v -o authorization ./cmd/controller

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/authorization /authorization

# Set the entry point
CMD ["/authorization"]
