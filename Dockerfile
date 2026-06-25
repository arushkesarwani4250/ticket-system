# Stage 1: Build stage
FROM golang:alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project source code
COPY . .

# Build the Go application statically linked (CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ticket-system cmd/server/main.go

# Stage 2: Final minimal image for execution
FROM alpine:3.19

# Add CA certificates (needed for secure outgoing requests like database TLS)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the compiled binary and Swagger docs from the builder stage
COPY --from=builder /app/ticket-system .
COPY --from=builder /app/api ./api

# Expose port 8080 as required by the brief contract
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./ticket-system"]
