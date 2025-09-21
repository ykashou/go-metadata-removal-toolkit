# Build Operations

This directory contains container build definitions for the go-metadata-removal-toolkit project.

## Container Files

- `development.Containerfile` - Development build with debugging tools
- `master.Containerfile` - Production build with distroless base

## Build Commands

### Development Build
```bash
# Build development container
podman build -f ops/build/development.Containerfile -t metadata-remover:dev .

# Run development container with mounted workspace
podman run --rm -it -v $(pwd):/workspace metadata-remover:dev
```

### Production Build
```bash
# Build production container
podman build -f ops/build/master.Containerfile -t metadata-remover:latest .

# Run production container to process files
podman run --rm -v /path/to/files:/data metadata-remover:latest -path /data -recursive
```

## Multi-Architecture Builds

```bash
# Build for multiple architectures
podman buildx build \
  --platform linux/amd64,linux/arm64 \
  -f ops/build/master.Containerfile \
  -t metadata-remover:latest \
  --push .
```

## Security Notes

- Production builds use distroless images for minimal attack surface
- Runs as non-root user in production
- No shell or package manager in production image
- All builds use specific version tags (no :latest)
- Binary is statically compiled for portability
