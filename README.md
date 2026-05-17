# Mini Docker

A learning project exploring core Linux container primitives in Go.

## Current State

The current implementation uses a re-exec flow to start an isolated `/bin/sh` process:

- The parent process (`runtime.Run()`) checks for an existing root filesystem and re-executes the current binary (`/proc/self/exe`) with new PID, UTS, and mount namespaces.
- The re-executed child process (`runtime.ContainerInit()`) isolates the mount namespace, configures an overlay filesystem, pivots the root (`pivot_root`), mounts `/proc`, and starts `/bin/sh`.

The process is attached to the caller's terminal through standard input, output, and error streams.

## Structure

- `main.go`: program entry point; branches between parent and child execution based on the `CONTAINER_INIT` environment variable.
- `runtime/run.go`: handles both parent (namespace setup) and child (filesystem isolation, shell execution) lifecycles.
- `runtime/rootfs.go`: handles the `pivot_root` functionality for the new filesystem root.
- `runtime/utils.go`: utility functions for checking and creating directories.
- `cli/run.go`: package stub reserved for future CLI work

## Planned Direction

The broader project is intended to grow toward:

- more namespace isolation (e.g., network, user)
- cgroup-based resource limits
- dynamic image handling, unpacking, and cleanup
- command parsing and additional runtime lifecycle features

*Note: This is an educational tool, not for production.*
