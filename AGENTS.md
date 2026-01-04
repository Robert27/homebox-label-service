# Repository Guidelines

## Project Structure & Module Organization
- `src/`: Go service source (entrypoint in `main.go`, HTTP handling in `handler.go`, rendering in `render.go`/`png.go`).
- `vendor/`: vendored Go dependencies for reproducible builds.
- `flake.nix`/`flake.lock`: Nix build definition and pinned inputs (also used to build the Docker image in CI).
- `README.md`: runtime behavior, query parameters, and environment variables.

## Build, Test, and Development Commands
- `go run ./src`: run the label service locally (defaults to `:8080`).
- `go build ./src`: build the service binary.
- `go test ./...`: run tests (no tests currently; use this when adding new ones).
- `gofmt -w src/*.go`: apply Go formatting before committing.
- `nix build .#homebox-label-service`: build via Nix (matches CI packaging).
- `nix build .#packages.x86_64-linux.dockerImage`: build the local Docker image output.

## Coding Style & Naming Conventions
- Follow standard Go style; rely on `gofmt` for formatting (tabs for indentation).
- Keep files in `package main` and name new files in lower-case `*.go` (e.g., `metrics.go`).
- Use Go naming: exported identifiers in `CamelCase`, unexported in `camelCase`.

## Testing Guidelines
- Tests live alongside code in `src/` and should be named `*_test.go`.
- Use the standard `testing` package and keep tests focused on rendering edge cases and parameter parsing.
- Run `go test ./...` locally before opening a PR.

## Commit & Pull Request Guidelines
- Commit messages follow Conventional Commits (e.g., `fix: ...`, `chore: ...`, `refactor: ...`).
- PRs should include a concise summary, test command(s) run, and any relevant issue links.
- If output rendering changes, attach a sample PNG (e.g., from the `curl` example in `README.md`).

## Configuration & Environment
- `PORT` controls the listen port (default `8080`).
- `HBOX_LABEL_MAKER_LABEL_SERVICE_TIMEOUT` sets request timeouts (e.g., `30s`).
- `HBOX_WEB_MAX_UPLOAD_SIZE` caps response size in bytes.
- `HBOX_LABEL_MAKER_LABEL_SERVICE_URL` is the URL Homebox should call.
