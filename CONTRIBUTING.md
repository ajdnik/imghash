# Contributing to imghash

Thank you for considering a contribution to imghash! This document explains how
to set up a development environment, run the checks that CI enforces, and submit
a pull request.

## Prerequisites

| Tool | Minimum version | Install |
|------|----------------|---------|
| Go | 1.16+ | [go.dev/dl](https://go.dev/dl/) |
| golangci-lint | v2.10.1 | Installed automatically by `make lint` |
| GNU Make | any | Pre-installed on most systems |

## Getting Started

```bash
git clone https://github.com/ajdnik/imghash.git
cd imghash
go mod download
```

## Development Workflow

The Makefile provides every command you need:

```bash
make fmt        # Format code with go fmt
make vet        # Run go vet
make lint       # Run golangci-lint (installs it if missing)
make test       # Run the test suite
make coverage   # Run tests with coverage and print a summary
make all        # fmt + vet + lint + test (recommended before pushing)
```

Run `make all` before opening a pull request — it mirrors the CI checks that
must pass.

## Submitting Changes

1. **Fork** the repository and create a feature branch from `main`:

   ```bash
   git checkout -b feature/my-change main
   ```

2. **Make your changes.** Keep each commit focused on a single logical change.

3. **Follow Conventional Commits.** This project uses
   [release-please](https://github.com/googleapis/release-please) to automate
   releases, so commit messages must follow the
   [Conventional Commits](https://www.conventionalcommits.org/) specification:

   ```
   feat: add Whirlpool hash algorithm
   fix: correct median calculation for even-length slices
   docs: clarify interpolation method defaults
   test: add edge-case coverage for empty images
   chore: bump golangci-lint to v2.11.0
   ```

   A `feat` commit triggers a minor version bump; a `fix` commit triggers a
   patch bump. Add `BREAKING CHANGE:` in the commit body (or append `!` after
   the type) to trigger a major bump.

4. **Open a pull request** against `main`. The following CI checks are required
   to pass:

   - **Go Vet** — `go vet ./...`
   - **Lint** — `golangci-lint run ./...`
   - **Test & Coverage** — `go test` with an 80 % minimum coverage threshold

5. A maintainer will review your PR. Please be responsive to feedback — small
   follow-up commits are fine; they get squash-merged.

## Reporting Bugs & Requesting Features

Open an [issue](https://github.com/ajdnik/imghash/issues) with a clear title
and as much context as possible. For bugs, include the Go version, OS, a
minimal reproducer, and the expected vs. actual behavior.

## Code Guidelines

- **No CGo.** The library is pure Go by design.
- **Test against OpenCV.** Hash algorithm implementations are ported from
  [OpenCV Contrib](https://github.com/opencv/opencv_contrib). When adding or
  modifying an algorithm, verify results against the OpenCV implementation.
- **Exported API surface.** Think twice before adding to the public API. Prefer
  functional options (`With*` pattern) over new struct fields.
- **Test coverage.** New code should be accompanied by tests. The CI enforces an
  80 % minimum coverage threshold — aim to stay above it.

## License

By contributing you agree that your contributions will be licensed under the
[MIT License](LICENSE).
