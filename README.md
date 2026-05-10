# Mini Docker

A learning project implementing core Linux container primitives in Go.

## Core Isolation
- **Namespaces**: PID, Mount, UTS, Network, IPC.
- **Cgroups v2**: Resource limits (CPU, Memory).
- **Filesystem**: OverlayFS + `pivot_root`.
- **Go Pattern**: Re-exec for namespace cloning.

## Structure
- `cli/`: Entry point.
- `runtime/`: Lifecycle management.
- `cgroups/`: Resource constraints.
- `network/`: veth networking.
- `image/`: Rootfs handling.

*Note: This is an educational tool, not for production.*
