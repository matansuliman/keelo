# Contributing to Keelo

First off, thank you for considering contributing to Keelo! 

## Project Architecture
Keelo is built in Go and acts as a package manager and template renderer for Docker Compose.
Before contributing, please read the `README.md` files located in the various subdirectories (`cmd/`, `internal/`, `pkg/`) to understand how the codebase is structured.

### Core Workflow
1. **Config Parsing**: Reads `project.yaml`.
2. **Module Loading**: Fetches definitions from `modules/` (local or remote git/oci registries).
3. **Template Rendering**: Merges inputs provided into Docker Compose templates using Go's `text/template`.
4. **Composition**: Merges rendered YAML fragments into a final `docker-compose.generated.yaml`.

## Development Setup
1. Ensure you have Go 1.24+ installed.
2. Fork and clone the repository.
3. Run `go mod download`.
4. Build the CLI with `go build -o keelo ./cmd/tool/`.

## Pull Request Process
1. **Branch off `dev`**: The `main` branch is strictly for stable releases. All new features and fixes should be branched off `dev`.
2. **Write Tests**: Ensure any new logic is covered by tests in the respective package. Run `go test ./...` to verify.
3. **Keep it Clean**: Write clear commit messages.
4. **Target `dev`**: Open your Pull Request against the `dev` branch. Our CI pipeline will automatically run the tests.

## Code of Conduct
Please communicate respectfully and constructively in Issues and Pull Requests. Let's build something great together!
