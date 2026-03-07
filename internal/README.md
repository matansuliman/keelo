# `internal/` - Private Application and Library Code

This directory contains the core business logic of Keelo.
In Go, the `internal/` directory is special: packages defined here can **only** be imported by code inside this repository. This prevents external tools from depending on Keelo's unstable internal APIs.

## Sub-Packages and Files

### `cli/`
Command-line interface definitions using Cobra/Viper. Wires commands to their functions.
*   **`root.go`**: Defines the base `keelo` command and global flags.
*   **`init.go`**: Scaffolds a new project by generating a base `project.yaml`.
*   **`list_modules.go`**: Prints available local/remote modules.
*   **`validate.go`**: Checks user inputs against module schemas without rendering.
*   **`render.go`**: Loads, validates, renders, and merges the project into a Compose file.
*   **`get.go`**: Triggers remote module downloading and caching.
*   **`diff.go`**: Renders the project in-memory and compares it to the existing generated file.
*   **`exec_commands.go`**: Wrappers for `docker compose up`, `down`, and `logs`.
*   **`cli_e2e_test.go`**: End-to-end tests validating the full CLI execution flow.

### `compose/`
Logic related to the final Docker Compose output.
*   **`writer.go`**: Handles taking the final merged byte slice and safely writing it to `docker-compose.generated.yaml` with a warning header.
*   **`writer_test.go`**: Unit tests for the writer.

### `config/`
Logic for parsing the user's project definition.
*   **`loader.go`**: Reads `project.yaml` and unmarshals it into Go structures.
*   **`loader_test.go`**: Unit tests for config loading.
*   **`lock.go`**: Automatically generates and loads `keelo.lock` to track remote module versions.

### `exec/`
System execution handlers.
*   **`exec.go`**: Safely wraps `os/exec` to run system commands (like `docker compose`) and pipe identical standard output/error to the terminal.
*   **`exec_test.go`**: Unit tests for system command execution.

### `merger/`
The YAML merging engine.
*   **`merger.go`**: Intelligently takes disparate YAML fragments (e.g., multiple modules outputting `services:`) and merges them into a single coherent schema while detecting top-level key conflicts.
*   **`merger_test.go`**: Extensive unit tests validating complex YAML merge boundaries.

### `modules/`
Module discovery and fetching.
*   **`loader.go`**: Searches the `modules/` directory for `module.yaml` files and parses their schemas into memory.
*   **`loader_test.go`**: Unit tests for local module parsing.
*   **`getter.go`**: Leverages HashiCorp's `go-getter` to securely download, unpack, and cache remote module repositories.

### `renderer/`
Template manipulation.
*   **`renderer.go`**: Uses standard Go `text/template` on raw `compose.yaml.tmpl` files, injecting user-provided input values securely into the templates.
*   **`renderer_test.go`**: Unit tests for template behavior and syntax.

### `validator/`
Input validation.
*   **`validator.go`**: Ensures that every value a user provides in `project.yaml` matches the required keys and type definitions (string, int, etc.) demanded by the corresponding `module.yaml`.
*   **`validator_test.go`**: Unit tests for schema enforcement constraints.
