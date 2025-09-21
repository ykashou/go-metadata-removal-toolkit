<div align="center">

# Go Metadata Removal Toolkit

A fast and efficient command-line tool for removing metadata from various file formats

[![Build Status](https://github.com/ykashou/go-metadata-removal-toolkit/workflows/CI/badge.svg)](https://github.com/ykashou/go-metadata-removal-toolkit/actions)
[![Test Coverage](https://codecov.io/gh/ykashou/go-metadata-removal-toolkit/branch/main/graph/badge.svg)](https://codecov.io/gh/ykashou/go-metadata-removal-toolkit)
[![License: ACE](https://img.shields.io/badge/License-ACE-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)

[Documentation](./docs) • [Report Bug](https://github.com/ykashou/go-metadata-removal-toolkit/issues) • [Request Feature](https://github.com/ykashou/go-metadata-removal-toolkit/issues)

</div>

## 🚀 Features

- ✅ **Multi-format Support**: Process images, PDFs, and documents
- ✅ **Batch Processing**: Handle multiple files and directories
- ✅ **Recursive Scanning**: Process entire directory trees
- ✅ **Preview Mode**: Safe dry-run before actual metadata removal
- ✅ **Detailed Statistics**: Track processing results and performance
- ✅ **Colored Output**: Beautiful terminal interface with progress indicators
- ✅ **Comprehensive Testing**: Unit tests with 1:1 source-to-test ratio
- ✅ **Containerized**: Run with Podman for consistent environments

## 🛠 Tech Stack

- **Language**: Go 1.19+
- **Logging**: Custom colored terminal output
- **File Processing**: Parallel processing with goroutines
- **Testing**: Go testing framework
- **Containerization**: Podman with distroless images
- **CI/CD**: GitHub Actions

## 📋 Prerequisites

- Go >= 1.19
- Podman >= 4.0 (optional, for containerized usage)

## 🚦 Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/ykashou/go-metadata-removal-toolkit.git
cd go-metadata-removal-toolkit

# Build the binary
go build -o metadata-remover

# Run the tool
./metadata-remover --help
```

### Basic Usage

```bash
# Remove metadata from a single file
./metadata-remover -path /path/to/file.jpg

# Process directory recursively
./metadata-remover -path /path/to/directory -recursive

# Preview mode (no changes made)
./metadata-remover -path /path/to/directory -preview

# Verbose output
./metadata-remover -path /path/to/directory -verbose
```

### Using Podman

```bash
# Build container image
podman build -t metadata-remover .

# Run containerized
podman run -v ./input:/data metadata-remover -path /data -recursive
```

## 📁 Project Structure

```
go-metadata-removal-toolkit/
├── src/                   # Source code directory
│   ├── logger/           # Logging utilities
│   ├── processor/        # File processing logic
│   │   ├── document.go   # Document metadata handler
│   │   ├── image.go      # Image metadata handler
│   │   └── pdf.go        # PDF metadata handler
│   ├── scanner/          # Directory scanning
│   ├── stats/            # Statistics collection
│   ├── utils/            # Utility functions
│   └── main.go           # Application entry point
├── ops/                  # Operations directory
│   ├── build/           # Build configurations
│   ├── release/         # Release management
│   ├── security/        # Security configurations
│   └── test/            # Test configurations
├── docs/                # Documentation
├── .github/             # GitHub Actions workflows
├── Containerfile        # Container build instructions
├── go.mod              # Go module definition
└── README.md           # This file
```

## 🧪 Testing

Maintaining a 1:1 source-to-test file ratio:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐳 Containerization

### Development

```bash
# Build development container
podman build -f ops/build/Containerfile.dev -t metadata-remover:dev .

# Run development container
podman run -it --rm -v $(pwd):/workspace metadata-remover:dev
```

### Production

```bash
# Build production container
podman build -f ops/build/Containerfile -t metadata-remover:latest .

# Run production container
podman run --rm -v /path/to/files:/data metadata-remover:latest -path /data
```

## 🔧 Configuration

### Command-line Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--path` | `-p` | Path to directory or file to process | `.` |
| `--recursive` | `-r` | Recursively process subdirectories | `false` |
| `--preview` | | Preview mode (no actual changes) | `false` |
| `--verbose` | `-v` | Verbose output | `false` |
| `--output` | | Output format (terminal, json) | `terminal` |
| `--version` | | Show version information | `false` |

## 📊 Repository Stats

![Repobeats](https://repobeats.axiom.co/api/embed/go-metadata-removal-toolkit.svg "Repobeats analytics image")

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ykashou/go-metadata-removal-toolkit&type=Date)](https://star-history.com/#ykashou/go-metadata-removal-toolkit&Date)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Development Guidelines

- Follow 1:1 source:test ratio
- Use conventional commits
- Ensure all tests pass
- Update documentation
- Use `gofmt` for code formatting
- Run `go vet` before committing

## 📄 License

This project is licensed under the ACE License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Contributors](https://github.com/ykashou/go-metadata-removal-toolkit/graphs/contributors) who helped build this project
- Go community for excellent tooling and libraries

---

<div align="center">
Made with ❤️ by <a href="https://github.com/ykashou">ykashou</a>
</div>
