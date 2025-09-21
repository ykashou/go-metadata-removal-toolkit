# Production Containerfile for go-metadata-removal-toolkit
FROM golang:1.19-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies (if go.sum exists)
RUN if [ -f go.sum ]; then go mod download; else go mod tidy; fi

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o metadata-remover ./src

# Production stage - distroless
FROM gcr.io/distroless/static:nonroot

# Copy binary from builder
COPY --from=builder /app/metadata-remover /usr/local/bin/metadata-remover

# Create data directory for processing
WORKDIR /data

# Use non-root user
USER nonroot:nonroot

# Run the application
ENTRYPOINT ["/usr/local/bin/metadata-remover"]
