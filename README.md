# Keelo

> A modular composition layer for Docker Compose. **Build infrastructure, not YAML.**

**Keelo** (formerly `tool` in development) is a CLI package manager and renderer for Docker Compose. It solves the "Giant YAML" problem by allowing you to define infrastructure components as **Reusable Modules** and stitch them together into a final deployment.

## Why Keelo?

Managing Docker Compose for multiple microservices or complex environments often leads to:
1.  **YAML Bloat**: Thousands of lines in a single `docker-compose.yaml`.
2.  **Duplication**: Copy-pasting the same Postgres/Redis configuration across 10 projects.
3.  **Human Error**: Forgetting to update a volume name or container port in one of many static files.

**Keelo fixes this by treating infrastructure like code modules.** You define a component (like a Database or an API) *once*, and reuse it across *any* number of projects with custom inputs.

---

## Core Concepts

### 1. The Module (The Building Block) 🧱
A folder containing a schema (`module.yaml`) and a template (`compose.yaml.tmpl`). 
- **Purpose**: Defines *how* a specific service should run.
- **Reusability**: You create it once in your `modules/` library and use it forever.

### 2. The Project (The Assembly) 🏗️
A `project.yaml` file that lists which modules you want and what values to give them.
- **Purpose**: Defines *what* your specific stack looks like.
- **Outcome**: `keelo render` compiles everything into a single, optimized `docker-compose.generated.yaml`.

---

## Project Structure
To help contributors navigate the repository, here is a breakdown of the core directories and the essential files living in the root:

### Core Directories
*   [**`cmd/`**](cmd/README.md): CLI entry points (the `keelo` binary). *See the directory's README for file details.*
*   [**`internal/`**](internal/README.md): Core private business logic (CLI routing, Config, Merging, Modules, Templates). *See the directory's README for file details.*
*   [**`pkg/`**](pkg/README.md): Public shared types (`ProjectConfig`, `ModuleDefinition`). *See the directory's README for file details.*
*   [**`examples/`**](examples/README.md): Demo projects showcasing Keelo's capabilities.

### Root Files
*   **`.github/workflows/`**: Contains our CI/CD pipelines. Includes `release.yml` (for cutting multi-platform binaries) and `test.yml` (for running PR tests).
*   **`.gitignore`**: Excludes raw binaries, temporary logs, and dynamic IDE files.
*   **`.goreleaser.yaml`**: The configuration file that tells GoReleaser how to build, package, and distribute the CLI across Homebrew, Scoop, and standard archives.
*   **`CONTRIBUTING.md`**: Community guidelines for contributing, writing tests, and opening PRs.
*   **`README.md`**: This file! The ultimate guide to Keelo.
*   **`go.mod` / `go.sum`**: Go's package tracking files. Defines the external dependencies required to run Keelo (e.g., Cobra, Viper, go-getter).

---

## Installation

### macOS & Linux (via Homebrew)
```bash
brew tap matansuliman/tap
brew install keelo
```

### Windows (via Scoop)
```powershell
scoop bucket add matansuliman https://github.com/matansuliman/homebrew-tap
scoop install keelo
```

### Direct Download (All Platforms)
Download the latest pre-compiled binary for Windows, macOS, or Linux from the [Releases](https://github.com/matansuliman/keelo/releases) page. Extract the `.tar.gz` or `.zip` file and place the `keelo` executable in your system's PATH.

### From Source
Ensure you have [Go](https://go.dev/doc/install) installed (1.24+ recommended).

```bash
go install github.com/matansuliman/keelo/cmd/tool@latest
```
*(This will automatically download, compile, and place the `keelo` binary in your `GOPATH/bin` directory).*

---

## Quick Start: From Zero to Deployment

### 1. Create your Library
Imagine `modules/` is your personal app store of infrastructure. Let's create a **Postgres** module.

**File: `modules/postgres/module.yaml`** (The Interface)
```yaml
name: postgres
description: "Reusable Postgres DB"
inputs:
  VERSION:
    type: string
    default: "15-alpine"
  DB_NAME:
    type: string
    required: true
```

**File: `modules/postgres/compose.yaml.tmpl`** (The Blueprint)
```yaml
services:
  db-{{ .ProjectName }}:
    image: postgres:{{ .Values.VERSION }}
    environment:
      POSTGRES_DB: {{ .Values.DB_NAME }}
    volumes:
      - pgdata-{{ .ProjectName }}:/var/lib/postgresql/data
volumes:
  pgdata-{{ .ProjectName }}:
```

### 2. Assembly your Project
Now, instead of writing 15 lines of Postgres YAML, you just "call" your module in `project.yaml`.

Run the interactive `init` command to bootstrap your configuration:
```bash
keelo init my-awesome-app
```
*The CLI will ask you for a starting template (e.g. Basic Web Service, Full Stack). Follow the prompts to automatically generate your `project.yaml`.*

Modify **`project.yaml`**:
```yaml
project: my-awesome-app
modules:
  - name: postgres
    values:
      DB_NAME: "user_service_db"
```

### 3. Generate and Run
Keelo handles the merging, input validation, and rendering.

```bash
# 🔍 List your library
keelo list-modules

# ✅ Validate inputs
keelo validate

# 🚀 Render and Start
keelo render
keelo up
```

---

## How it Saves You Time
*   **Zero Boilerplate**: Want a Redis instance? Just add 3 lines to `project.yaml` instead of 20 lines of Compose.
*   **Type Safety**: Keelo validates that you provided all `required` inputs and that they are the right type (string, int, etc).
*   **Dynamic Names**: Module templates use `{{ .ProjectName }}` to ensure that if you run two different projects, their container and volume names never collide.
*   **Centralized Updates**: Update the `postgres` module template once, and every project using it gets the update on the next `render`.

---

## License
MIT
