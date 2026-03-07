# Example Projects

This directory contains functional examples of how to compose infrastructure using Keelo.
These are great starting points if you are learning how to use the tool.

## Available Examples

*   **`01-basic-postgres`**: The simplest possible example. It loads a single local Postgres module, demonstrates default inputs, and creates a basic database instance.
    *   `project.yaml`: The project definition listing the `postgres` module.
    *   `modules/`: A local library containing the reusable postgres definition.
*   **`02-full-stack-web`**: A foundational multi-tier architecture. It composes a frontend API service, a Redis caching layer, and a Postgres database, showing how multiple modules work together in a single `project.yaml`.
    *   `project.yaml`: Lists three separate modules (`api`, `db`, `cache`).
    *   `modules/`: Contains the three local module definitions.
*   **`03-secrets-env`**: Demonstrates Keelo's secure environment variable injection. It shows how you can reference secrets in `project.yaml` without hardcoding them into your repository.
    *   `project.yaml`: Uses `${ENV_VAR}` syntax.
    *   `.env.example`: Shows how to pass environment variables to the CLI.
*   **`04-remote-modules`**: Shows off Keelo's package management capabilities by fetching modules directly from a remote Git repository instead of a local folder.
    *   `project.yaml`: References Git URLs instead of local file paths.
    *   `.keelo/`: (Auto-generated) Local cache where remote modules are downloaded and stored.
