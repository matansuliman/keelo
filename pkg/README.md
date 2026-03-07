# `pkg/` - Public Library Code

This directory contains packages that are safe to be used by external Go applications.
Unlike `internal/`, the code here constitutes a public API.

## `pkg/types`
This package contains the core data structures used throughout the Keelo application. By keeping these structs here, different internal subsystems (like the config parser, the template renderer, and the CLI) can share a common vocabulary without creating circular dependencies.

### Key Types
*   **`ProjectConfig`**: Represents the parsed state of a user's `project.yaml`.
*   **`ModuleDefinition`**: Represents the schema loaded from a module's `module.yaml`.
*   **`RenderedModule`**: Represents the intermediate state of a module after its template has been executed, but before it is merged into the final Compose file.
