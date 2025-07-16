# goalias

[![Test](https://github.com/jackchuka/goalias/workflows/Test/badge.svg)](https://github.com/jackchuka/goalias/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackchuka/goalias)](https://goreportcard.com/report/github.com/jackchuka/goalias)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A fast and reliable Go CLI tool for managing import aliases consistently across your entire codebase using the Language Server Protocol (LSP).

## Features

- 🚀 **Fast**: Uses persistent LSP connections for processing multiple files
- 🔍 **Accurate**: Leverages `gopls` for precise Go code analysis and refactoring
- 📦 **Comprehensive**: Scans entire modules or specific package patterns
- 🛡️ **Safe**: Skips generated files and provides clear feedback
- 🎯 **Consistent**: Ensures uniform import aliases across your codebase
- 📊 **Informative**: Shows detailed location information for all import usages

## Installation

```bash
# Install goalias
go install github.com/jackchuka/goalias/cmd/goalias@latest

# Install gopls (required dependency)
go install golang.org/x/tools/gopls@latest
```

Both `goalias` and `gopls` must be available in your `$PATH`.

## Quick Start

```bash
# List all current aliases for a package
goalias list -p github.com/example/mypackage

# Set a consistent alias across your entire module
goalias set -p github.com/example/mypackage -a mypkg

# Set alias for specific package patterns
goalias set -p github.com/example/mypackage -a mypkg ./cmd/... ./internal/...
```

## Usage

### Set Import Alias

Update import aliases consistently across your codebase:

```bash
# Set alias for entire module
goalias set -p github.com/example/mypackage -a mypkg

# Set alias for specific package patterns
goalias set -p github.com/example/mypackage -a mypkg ./cmd/...

# Multiple patterns
goalias set -p github.com/example/mypackage -a mypkg ./cmd/... ./internal/...
```

### List Import Aliases

Display all current aliases and their locations:

```bash
# List aliases for entire module
goalias list -p github.com/example/mypackage

# List aliases for specific package patterns
goalias list -p github.com/example/mypackage ./cmd/...
```

**Example output:**

```
LOCATION                  ALIAS
--------                  -----
handler/foo.go:4         myutils
handler/bar.go:6         myutils
cmd/server/main.go:8     mypackage
```

## Commands

### `goalias set`

Sets or updates import aliases across specified packages.

```bash
goalias set --package <importPath> --alias <name> [patterns...]
```

**Required Flags:**

- `--package`, `-p`: Full import path to manage
- `--alias`, `-a`: Desired alias identifier

**Optional Arguments:**

- `patterns`: Go package patterns (defaults to `./...`)

**Examples:**

```bash
# Set alias for entire module
goalias set -p github.com/gin-gonic/gin -a gin

# Set alias for specific directories
goalias set -p github.com/stretchr/testify/assert -a tassert ./tests/...

# Multiple patterns
goalias set -p github.com/pkg/errors -a pkg_errros ./cmd/... ./internal/...
```

### `goalias list`

Lists all import aliases and their locations.

```bash
goalias list --package <importPath> [patterns...]
```

**Required Flags:**

- `--package`, `-p`: Full import path to search for

**Optional Arguments:**

- `patterns`: Go package patterns (defaults to `./...`)

**Examples:**

```bash
# List all usages in module
goalias list -p github.com/gin-gonic/gin

# List usages in specific directories
goalias list -p github.com/stretchr/testify/assert ./tests/...
```

## How It Works

1. **Package Discovery**: Uses `go list` to find Go packages matching your patterns
2. **AST Parsing**: Parses Go source files to locate import declarations
3. **Smart Filtering**: Automatically skips generated files (containing `Code generated ... DO NOT EDIT`)
4. **LSP Integration**: Uses persistent `gopls` connections for accurate code analysis and refactoring
5. **Consistent Updates**: Applies import alias changes atomically across all matching files

## Performance

goalias is designed for speed and efficiency:

- **Persistent LSP Connections**: Reuses `gopls` sessions instead of spawning new processes
- **Batch Operations**: Processes multiple files in a single LSP session
- **Smart Caching**: Leverages `gopls` internal caching for faster analysis

## Requirements

- Go 1.22 or later
- `gopls` (Go Language Server)
- Git (for repository operations)

## Development

### Building from Source

```bash
git clone https://github.com/jackchuka/goalias.git
cd goalias
go build -o goalias ./cmd/goalias
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

### Development Setup

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run tests (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style

- Follow standard Go formatting (`golangci-lint`)
- Add tests for new functionality
- Update documentation as needed
- Follow existing code patterns and conventions

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [gopls](https://pkg.go.dev/golang.org/x/tools/gopls) for Go language server capabilities
- Inspired by the need for consistent import management in large Go codebases

## Support

If you encounter any issues or have questions:

1. Check the existing [issues](https://github.com/jackchuka/goalias/issues)
2. Create a new issue with detailed information
3. Include your Go version, OS, and example code when reporting bugs
