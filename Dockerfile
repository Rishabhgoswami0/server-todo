# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the workspace file to the builder root
COPY go.work go.work.sum* ./

# Copy the shared module
COPY shared-go/ ./shared-go/

# Copy the server-todo module
COPY server-todo/ ./server-todo/

# Navigate to the server-todo module directory
WORKDIR /app/server-todo

# Download dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

# Install necessary certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server-todo/main .

# Expose the API port
EXPOSE 3002

# Command to run the application
CMD ["./main"]
