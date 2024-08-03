# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install git to support go modules
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /myredis

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .


RUN ls -al
RUN go list -m all
# Build the Go binary
RUN go build -o main .

# Stage 2: Run the Go binary
FROM alpine:latest

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /myredis/main .

# Expose port 8080 to the outside world
EXPOSE 7379

# Command to run the executable
CMD ["./main"]