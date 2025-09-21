# Release Operations

This directory contains release management scripts and configurations for the go-metadata-removal-toolkit project.

## Files

- `release-config.yaml` - Release configuration parameters
- `release.sh` - Main release automation script
- `version-bump.sh` - Version increment utility

## Release Process

### 1. Prepare Release
```bash
# Check current version
./ops/release/version-bump.sh current

# Bump version (major, minor, or patch)
./ops/release/version-bump.sh patch
```

### 2. Create Release
```bash
# Run release script
./ops/release/release.sh

# Or specify version explicitly
./ops/release/release.sh v1.0.1
```

### 3. Release Steps
The release script automatically:
1. Validates version format
2. Runs all tests
3. Builds binaries for multiple platforms
4. Creates container images
5. Generates changelog
6. Tags the release
7. Pushes to registry

## Version Scheme
We follow semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Incompatible API changes
- MINOR: New functionality (backward compatible)
- PATCH: Bug fixes (backward compatible)

## Platform Targets
- linux/amd64
- linux/arm64
- darwin/amd64
- darwin/arm64
- windows/amd64

## Release Artifacts
- Binary executables per platform
- Container images (multi-arch)
- Source code archives
- Checksums (SHA256)
- Release notes
