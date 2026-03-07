# Example 03: Secrets and Environment Variables

This example demonstrates how to use Keelo's template functions to inject secrets and dynamic configuration from environment variables.

## Features
- **`env` function**: Pulls values directly from the host environment.
- **`default` filter**: Provides a fallback value if an environment variable or input is missing.
- **Sensitive data**: Marking inputs as `sensitive: true` in `module.yaml`.

## Files
- `project.yaml`: Sets the base environment and a static token.
- `modules/app/`: A module that uses both `env` and `default`.
- `.env.example`: Shows which environment variables can be set to customize the output.
- `generated/docker-compose.yaml`: The resulting Compose file with injected secrets.

## How to run
1. Set the environment variables:
   ```bash
   export EXTERNAL_API_KEY="prod-secret-999"
   export DEBUG_MODE="true"
   ```
2. Render the project:
   ```bash
   keelo render -c examples/03-secrets-env/project.yaml -o examples/03-secrets-env/generated/docker-compose.yaml
   ```

Check the `generated/docker-compose.yaml` to see how the secrets were injected.
