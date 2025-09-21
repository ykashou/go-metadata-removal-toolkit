# Development Containerfile for go-metadata-removal-toolkit
FROM golang:1.19-bullseye

# Install development tools
RUN apt-get update && apt-get install -y \
    git \
    vim \
    curl \
    wget \
    jq \
    make \
    && rm -rf /var/lib/apt/lists/*

# Install debugging tools
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set working directory
WORKDIR /workspace

# Copy go mod files
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

# Build with debug symbols
RUN go build -gcflags="all=-N -l" -o metadata-remover ./src

# Set environment for development
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Default command for development
CMD ["/bin/bash"]
