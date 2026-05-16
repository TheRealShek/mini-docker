# Mini Docker

A learning project exploring core Linux container primitives in Go.

## Current State

The current implementation uses a re-exec flow to start an isolated `/bin/sh` process:

- The parent process (`runtime.Run()`) re-executes the current binary (`/proc/self/exe`) with new PID, UTS, and mount namespaces.
- The re-executed child process (`runtime.ContainerInit()`) detects the `CONTAINER_INIT` environment variable and starts `/bin/sh`.

The process is attached to the caller's terminal through standard input, output, and error streams.

## Structure

- `main.go`: program entry point; branches between parent and child execution based on the `CONTAINER_INIT` environment variable.
- `runtime/run.go`: handles both parent (namespace setup) and child (shell execution) lifecycles.
- `cli/run.go`: package stub reserved for future CLI work

## Planned Direction

The broader project is intended to grow toward:

- more namespace isolation
- cgroup-based resource limits
- filesystem isolation with layered root filesystems
- command parsing and additional runtime lifecycle features

*Note: This is an educational tool, not for production.*
