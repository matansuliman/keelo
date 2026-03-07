# `internal/` - Private Application and Library Code

This directory contains the core business logic of Keelo.
In Go, the `internal/` directory is special: packages defined here can **only** be imported by code inside this repository. This prevents external tools from depending on Keelo's unstable internal APIs.

## Sub-Packages
*   `cli/`: Command-line interface definitions using Cobra/Viper. Wires commands like `keelo up` to their handler functions.
*   `compose/`: Logic related to interacting with and executing Docker Compose commands (wrappers for `os/exec`) and writing output files.
*   `config/`: Logic for parsing `project.yaml` and mapping it to the internal configuration structures.
*   `merger/`: The engine that takes multiple disparate YAML fragments and intelligently merges them (e.g., combining `services:`, `networks:`, `volumes:` maps) while detecting naming conflicts.
*   `modules/`: Logic for module discovery. Handles loading local modules from the file system and fetching remote modules (via Git/HTTP using `go-getter`).
*   `renderer/`: Uses Go's `text/template` engine to inject configured inputs into the raw `compose.yaml.tmpl` templates defined by modules.
*   `validator/`: Ensures that user inputs provided in `project.yaml` match the expected schema defined in the module's `module.yaml`.
