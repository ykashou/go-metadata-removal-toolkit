#!/bin/bash

# Release script for go-metadata-removal-toolkit
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="go-metadata-removal-toolkit"
BINARY_NAME="metadata-remover"
VERSION="${1:-}"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Validate version
if [ -z "$VERSION" ]; then
    log_error "Version not specified. Usage: $0 <version>"
fi

if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    log_error "Invalid version format. Expected: v1.0.0"
fi

log_info "Starting release for $PROJECT_NAME $VERSION"

# Run tests
log_info "Running tests..."
go test ./src/... || log_error "Tests failed"

# Build binaries for multiple platforms
log_info "Building binaries..."
mkdir -p dist

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    output="dist/${BINARY_NAME}-${VERSION}-${os}-${arch}"
    
    if [ "$os" = "windows" ]; then
        output="${output}.exe"
    fi
    
    log_info "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch go build -ldflags="-w -s -X main.appVersion=${VERSION}" \
        -o "$output" ./src || log_warning "Failed to build for $os/$arch"
done

# Generate checksums
log_info "Generating checksums..."
cd dist
sha256sum * > checksums-${VERSION}.sha256
cd ..

# Build container images
log_info "Building container images..."
podman build -f ops/build/master.Containerfile \
    -t ${PROJECT_NAME}:${VERSION} \
    -t ${PROJECT_NAME}:latest .

# Tag release in git
log_info "Tagging release..."
git tag -a "${VERSION}" -m "Release ${VERSION}"

log_info "Release $VERSION completed successfully!"
log_info "Next steps:"
echo "  1. Push tags: git push origin ${VERSION}"
echo "  2. Push container: podman push ${PROJECT_NAME}:${VERSION}"
echo "  3. Upload binaries from dist/ to GitHub releases"
