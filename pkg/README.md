# `pkg/` - Public Library Code

This directory contains packages that are safe to be used by external Go applications.
Unlike `internal/`, the code here constitutes a public API.

## `pkg/types`
This package contains the core data structures used throughout the Keelo application. By keeping these structs here, different internal subsystems (like the config parser, the template renderer, and the CLI) can share a common vocabulary without creating circular dependencies.

### Files inside `pkg/types`:
*   **`config.go`**: Defines `ProjectConfig` and related structs, representing the parsed state of a user's `project.yaml`.
*   **`lock.go`**: Defines the data model for the `keelo.lock` file, tracking dependencies.
*   **`module.go`**: Defines `ModuleDefinition` (the schema loaded from `module.yaml`) and `RenderedModule` (the intermediate state generated after template execution).
*   **`render.go`**: Defines the data structures passed directly into the Go `text/template` engine (e.g., `TemplateData`), controlling exactly what variables are exposed to template authors.
