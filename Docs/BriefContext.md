# Brief Context

## Project Essence

**Mini Docker** is a learning-oriented implementation of Linux container primitives built in Go. It aims to demonstrate how containers work from the ground up by directly interacting with kernel features.

- **Current Goal**: Build up a minimal container runtime incrementally, starting from isolated process execution without a background daemon.
- **Currently Implemented**:
  - `main.go` branches execution based on the `CONTAINER_INIT` environment variable.
  - `runtime.Run()` (parent) re-executes `/proc/self/exe` to create new PID, UTS, and mount namespaces.
  - `runtime.ContainerInit()` (child) detects the `CONTAINER_INIT` flag and starts `/bin/sh`.
  - Child stdio is connected to the invoking terminal.
- **Not Yet Implemented**:
  - CLI command parsing beyond the current direct entry point.
  - Cgroups, rootfs isolation, overlayfs, networking, image handling, and cleanup.

## Planned Direction

- `cli/`: Entry point and command parsing (run, ps, pull).
- `runtime/`: Container lifecycle management (creation, execution, cleanup).
- `cgroups/`: Resource limit application and management (v2 focused).
- `network/`: Virtual ethernet (veth) setup, IP assignment, and NAT rules.
- `image/`: Rootfs management, unpacking, and layer handling.


## Changelog

- 2026-05-16: Updated the brief context to reflect the newly implemented re-exec flow.
- 2026-05-16: Updated the brief context to distinguish the current implemented runtime from the planned architecture.

## Previous Essential State

- The document previously stated that `main.go` directly started `/bin/sh` and that the re-exec flow was not yet implemented.
- The document previously described the full target architecture as if PID, mount, UTS, network, IPC, cgroups v2, overlayfs, `pivot_root`, and re-exec were already part of the project state.
