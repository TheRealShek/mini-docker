# Brief Context

## Project Essence

**Mini Docker** is a learning-oriented implementation of Linux container primitives built in Go. It aims to demonstrate how containers work from the ground up by directly interacting with kernel features.

- **Goal**: Create an isolated process environment (container) without a background daemon.
- **Core Technologies**:
  - **Namespaces**: PID, Mount, UTS, Network, and IPC for isolation.
  - **Cgroups v2**: For resource limiting (CPU, Memory).
  - **Filesystem**: `overlayfs` for layered images and `pivot_root` for root filesystem isolation.
  - **Go Pattern**: Re-exec pattern to handle namespace cloning before the Go runtime fully initializes.

## Planned Structure

- `cli/`: Entry point and command parsing (run, ps, pull).
- `runtime/`: Container lifecycle management (creation, execution, cleanup).
- `cgroups/`: Resource limit application and management (v2 focused).
- `network/`: Virtual ethernet (veth) setup, IP assignment, and NAT rules.
- `image/`: Rootfs management, unpacking, and layer handling.

## AI/Agent Rules

1. **Read First**: Any AI or automated agent MUST read this file fully before analyzing the codebase, proposing changes, or running edits.
2. **Single Source of Context**: Treat this file as the primary short-form context for the project. If other docs conflict, call out the discrepancy.
3. **Detecting Drift**: If an agent finds the codebase and this document are out of sync (missing features, different behavior, or outdated notes), it MUST pause and notify the user.
4. **Ask Before Changing**: Before changing this document, the agent should ask the user for permission to update it.
5. **Update Procedure**: When updating, include a one-line changelog entry with date and short summary, and append a short backup block of the previous essential state.
6. **If Denied**: If the user denies updating the document, the agent should proceed only after documenting the discrepancy in a short note.
7. **Minimal Edits**: Keep edits factual, minimal, and human-readable.

---
