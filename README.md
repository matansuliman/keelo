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

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed (1.20+ recommended).

```bash
git clone https://github.com/matansuliman/keelo.git
cd keelo
go build -o keelo ./cmd/tool
# Add to your PATH
```

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

```bash
keelo init my-awesome-app
```

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

## Roadmap
*   [x] Phase 1-2: Core Config & Module Loading
*   [x] Phase 3-4: Template Rendering & YAML Merging
*   [x] Phase 5-6: Output Generation & Docker Compose Wrappers
*   [x] Phase 7-8: CLI Polish & E2E Testing
*   [x] Phase 9: Documentation 🚀

## License
MIT
