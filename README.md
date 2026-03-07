# Keelo

> A modular composition layer for Docker Compose

**Keelo** (formerly `tool` in development) is a CLI package manager and renderer for Docker Compose applications. It allows developers to declare infrastructure components dynamically as **Modules** rather than maintaining giant, static `docker-compose.yaml` files.

Think of it as Helm for Docker Compose. Instead of writing boilerplate YAML for every new project, you can reuse modular Compose definitions, supply configuration values, and `keelo` will stitch them together into a final, deployable `docker-compose.generated.yaml`.

## Features
- **Modular Architecture**: Break down infrastructure into reusable packages (e.g., `postgres`, `redis`, `api-gateway`).
- **Dynamic Configuration**: Pass inputs to modules using Go `text/template`.
- **Conflict Detection**: Automatically detects naming collisions in services, volumes, and networks across different modules.
- **Developer Experience**: Built-in commands to initialize projects, validate inputs, list available modules, and execute the generated Compose file.

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed (1.20+ recommended).

```bash
git clone https://github.com/matansuliman/keelo.git
cd keelo
go build -o keelo ./cmd/tool
# Move to your PATH, e.g., mv keelo /usr/local/bin/
```

## Quick Start

### 1. Initialize a Project

Create a new directory for your deployment and run:

```bash
mkdir my-deployment && cd my-deployment
keelo init my-project
```

This creates a `project.yaml` file:
```yaml
project: my-project
modules:
# ... add your modules here
```

### 2. Create a Module

Modules live in a `modules/` directory relative to your project. Let's create a minimal Redis module:

```bash
mkdir -p modules/redis
```

Create `modules/redis/module.yaml`:
```yaml
name: redis
version: "7.0"
description: "A simple Redis instance"
inputs:
  PORT:
    type: int
    default: 6379
```

Create `modules/redis/compose.yaml.tmpl`:
```yaml
services:
  {{ .ProjectName }}-redis:
    image: redis:{{ .Values.Version | default "latest" }}
    ports:
      - "{{ .Values.PORT }}:6379"
```

### 3. Configure the Project

Update your `project.yaml` to use the new module:

```yaml
project: demo
modules:
  - name: redis
    values:
      PORT: 6380
```

### 4. Validate and Render

Check that everything is configured correctly:

```bash
keelo validate
```

Render the final compose file:

```bash
keelo render
```
This generates a `docker-compose.generated.yaml` containing the merged results.

### 5. Execute Run

Stand up your infrastructure:

```bash
keelo up
# To view logs: keelo logs -f
# To tear down: keelo down
```

## Project Concept & Roadmap
This tool was built using a structured roadmap (Phases 0-9). Key capabilities implemented include:
* **Phase 1:** Project config parsing (`project.yaml`).
* **Phase 2:** Module definitions and input validation.
* **Phase 3 & 4:** Go template rendering and YAML merging.
* **Phase 5 & 6:** Final output generation and `os/exec` wrappers for Docker Compose.
* **Phase 7 & 8:** Improved developer experience, CLI commands, and test stabilization.
* **Phase 9:** Documentation.

## License
MIT License.
