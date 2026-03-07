# Example 04: Remote Modules

This example showcases Keelo's ability to fetch modules from remote Git repositories.

## Features
- **Remote Sourcing**: Using the `source` field in `project.yaml` to point to a Git repository.
- **Subdirectory Support**: Fetching a specific module from a subdirectory within a repository using the `//` syntax.
- **Caching**: Remote modules are automatically downloaded and cached in `.keelo/cache`.

## Files
- `project.yaml`: Points to a remote module source.
- `generated/docker-compose.yaml`: The resulting file rendered from the remote module.

## How to run
To run this example with a real remote module, you would update the `source` in `project.yaml` to a valid Git URL:
```yaml
modules:
  - name: my-mod
    source: "git::https://github.com/example/modules//my-mod"
```

Then run:
```bash
keelo render
```
Keelo will automatically clone the repository, cache it, and render the template.
