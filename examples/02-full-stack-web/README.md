# Example 02: Full-Stack Web Application

This example demonstrates a multi-module composition where different components (API, Database, Cache) are managed as separate modules and then combined into a single project.

## Features
- **Module Composition**: Three distinct modules working together.
- **Service Decoration**: Using the `{{ .ProjectName }}` variable to ensure unique service names.
- **Inter-Module Communication**: The API module is configured to connect to the DB and Redis services.

## Files
- `project.yaml`: Orchiestrates the modules and passes the necessary connection hosts.
- `modules/api/`: Frontend/API service with configurable port and dependencies.
- `modules/postgres/`: Database service.
- `modules/redis/`: Caching service.
- `generated/docker-compose.yaml`: The resulting merged Compose file.

## How to run
From the root of the Keelo repository, run:
```bash
keelo render -c examples/02-full-stack-web/project.yaml -o examples/02-full-stack-web/generated/docker-compose.yaml
```
