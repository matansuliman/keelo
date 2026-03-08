# Remote Keelo Modules Example

This example demonstrates how Keelo fetches, caches, and renders remote modules directly from a centralized remote repository, specifically `matansuliman/keelo-modules`.

In this example, we build a stack with a PostgreSQL database and a Redis cache without writing any Docker Compose boilerplate.

## 1. Inspect the Project

Look at `project.yaml`. Notice how the `source` attribute points directly to subdirectories within a GitHub repository using the `go-getter` syntax:
`git::https://github.com/matansuliman/keelo-modules//base-postgres`

## 2. Fetch and Render

Run the following command:

```bash
# This will download the remote modules to ~/.cache/keelo/modules/
# And generate docker-compose.generated.yaml
keelo render
```

If you re-run `keelo render`, Keelo will realize it already has the cached module versions locally and will render much faster!

## 3. Verify

Check the generated `docker-compose.generated.yaml`. You will see a complete, properly named, and fully configured compose file with volumes, ports, and injected environment variables based on the inputs provided in `project.yaml`.
