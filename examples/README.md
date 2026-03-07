# Keelo Examples 🚀

This directory contains various examples showcasing Keelo's capabilities, from simple setups to complex multi-module infrastructures.

## Directory Structure
- [01-basic-postgres](./01-basic-postgres/): Simple single-module project using PostgreSQL.
- [02-full-stack-web](./02-full-stack-web/): Multi-module composition with API, Redis, and Database.
- [03-secrets-env](./03-secrets-env/): Demonstrating environment variable injection and default fallbacks.
- [04-remote-modules](./04-remote-modules/): Using remote git repositories as module sources.

## How to use these examples
Each example is self-contained. You can run them by navigating to the example directory and running:
```bash
keelo render
```
The output will be generated in the `generated/` subdirectory of each example.
