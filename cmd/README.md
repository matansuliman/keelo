# `cmd/` - Command Line Applications

This directory contains the main entry points for the applications built from this repository.

## `cmd/tool/`
This is the root package for the `keelo` CLI binary. **Do not put business logic here.** This directory should remain as thin as possible.

### Files:
*   **`main.go`**: The absolute entry point of the Keelo application. Its only responsibility is to import `keelo/internal/cli` and call the Cobra `Execute()` function to trigger the command routing. If execution fails, it handles the standard `os.Exit(1)`.
