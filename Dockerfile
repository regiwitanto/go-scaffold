# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-scaffold ./cmd/scaffold

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /go-scaffold /app/
COPY --from=builder /app/templates /app/templates

# Expose port
EXPOSE 8081

# Run the application
CMD ["/app/go-scaffold"]
