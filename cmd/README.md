# `cmd/` - Command Line Applications

This directory contains the main entry points for the applications built from this repository.

## `cmd/tool`
This is the root package for the `keelo` CLI binary.
- It contains `main.go`, which simply initializes and executes the Cobra root command defined in `internal/cli`.
- **Do not put business logic here.** This directory should remain as thin as possible, acting only as the trigger for the internal package logic.
