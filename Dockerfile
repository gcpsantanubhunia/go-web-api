# --- STAGE 1: Build the binary ---
FROM golang:1.26.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency manifests first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build a statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /webapp .

# --- STAGE 2: Deploy minimal production container ---
FROM alpine:latest  

# Add a non-root user for production security
RUN adduser -D appuser
USER appuser

# Set working directory
WORKDIR /

# Copy only the compiled binary from the builder stage
COPY --from=builder /webapp /webapp

# Expose the application network port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/webapp"]
