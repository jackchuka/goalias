<div align="center">

# goalias

**One import path, one alias, across the whole codebase.**

Find every alias a Go package is imported under,<br>
then rename them all to one name with `gopls`, usages included.

[![Test](https://github.com/jackchuka/goalias/actions/workflows/test.yml/badge.svg)](https://github.com/jackchuka/goalias/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/jackchuka/goalias?sort=semver)](https://github.com/jackchuka/goalias/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

<img src="assets/demo.gif" alt="goalias listing three different aliases for one package, then setting them all to a single alias" width="800">

[Installation](#installation) · [Quick Start](#quick-start) · [Commands](#commands) · [How It Works](#how-it-works)

</div>

---

In a growing Go codebase the same package ends up imported as `pr` in one file, `p` in another and `pricingsvc` in a third. Renaming them by hand means editing the import line and every reference in each file. **goalias** lists every alias in use for an import path, then rewrites them to the one you choose through `gopls` rename, so each reference in the file changes with its import.

## Features

- **Audit** — `goalias list` shows every file that imports a package and the alias it uses
- **Rewrite** — `goalias set` renames the alias and all its usages, file by file, via `gopls`
- **Preview** — `--preview` prints the changes as a diff without writing anything
- **Scoped** — limit the scan to Go package patterns such as `./cmd/...`; defaults to `./...`
- **Safe defaults** — skips generated files (`Code generated ... DO NOT EDIT`) and files already using the target alias
- **Fast** — one persistent `gopls` session handles every file in a run

## Installation

### Go

```bash
go install github.com/jackchuka/goalias/cmd/goalias@latest
```

### Binary

Download a prebuilt binary for Linux, macOS or Windows from the [releases page](https://github.com/jackchuka/goalias/releases).

### Requirements

`goalias set` drives `gopls`, so it must be on your `$PATH`:

```bash
go install golang.org/x/tools/gopls@latest
```

`goalias list` only needs the `go` command. Building goalias from source requires Go 1.26 or later.

## Quick Start

Run these from the root of the Go module you want to clean up.

1. **See** which aliases are in use for a package:

   ```bash
   goalias list -p example.com/shop/pricing
   ```

   ```
   LOCATION                       ALIAS
   --------                       -----
   /src/shop/api/checkout.go:4    pr
   /src/shop/cmd/shop/main.go:6   p
   /src/shop/worker/invoice.go:6  pricingsvc
   ```

   Imports without an alias are reported under the package's default name (the last path element).

2. **Preview** the rename:

   ```bash
   goalias set -p example.com/shop/pricing -a pricing --preview
   ```

3. **Apply** it:

   ```bash
   goalias set -p example.com/shop/pricing -a pricing
   ```

   When the target alias matches the package name, `gopls` drops the explicit alias and leaves a plain import.

## Commands

| Command                                          | Description                                              |
| ------------------------------------------------ | -------------------------------------------------------- |
| `goalias list -p <path> [patterns...]`           | List every file importing `<path>` and the alias it uses |
| `goalias set -p <path> -a <alias> [patterns...]` | Rename the alias (and its usages) to `<alias>`           |
| `goalias version`                                | Print version, commit and build date                     |

`patterns` are Go package patterns passed to `go list`; they default to `./...`.

### Flags

| Flag            | Commands      | Description                                      |
| --------------- | ------------- | ------------------------------------------------ |
| `-p, --package` | `list`, `set` | Full import path to inspect or update (required) |
| `-a, --alias`   | `set`         | Desired alias identifier (required)              |
| `-n, --preview` | `set`         | Print a diff instead of writing files            |

### Examples

```bash
# Standardize across the whole module
goalias set -p github.com/pkg/errors -a pkgerrors

# Only in some packages
goalias set -p github.com/stretchr/testify/assert -a tassert ./internal/... ./cmd/...

# Audit one directory tree
goalias list -p github.com/gin-gonic/gin ./api/...
```

## How It Works

1. **Discover packages** — runs `go list -json` with your patterns and collects each package's Go files
2. **Parse imports** — parses each file's AST to find the import spec for the path, skipping generated files
3. **Rename via LSP** — opens one `gopls` session and issues a `textDocument/rename` on each import that doesn't already use the target alias
4. **Apply edits** — writes the returned workspace edits to disk, or prints them as a diff with `--preview`

Only non-test Go files (`GoFiles` from `go list`) are scanned; `_test.go` files are not changed.

## Development

```bash
git clone https://github.com/jackchuka/goalias.git
cd goalias
go build -o goalias ./cmd/goalias
go test ./...
```

The demo GIF is rendered with [VHS](https://github.com/charmbracelet/vhs): `vhs assets/demo.tape`.

## Contributing

Contributions are welcome. Please open an issue for bugs or feature ideas, and include your Go version, OS and a minimal example when reporting a bug. Pull requests should include tests and pass `golangci-lint`.

## License

[MIT](LICENSE)
