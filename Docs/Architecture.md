# Mini Docker

> This Docuemnt will be a Guiding Light as we Don't know how to approach a NEW Idea.

A learning project. Implements core container primitives using Linux
namespaces, cgroups, pivot_root, and overlayfs. Not production software.

---

## What It Does

Takes a command + rootfs path, creates an isolated process with:

- Own PID, mount, UTS, network namespaces
- Resource limits via cgroups
- Isolated filesystem via pivot_root + overlayfs
- Optional network via veth pair

---

## Components

### cli/

Parses subcommands: run, ps (later), pull (later).
Entry point. Calls into runtime.

### runtime/

Core logic. Owns the container lifecycle.

- container.go — Container struct, state (created/running/stopped)
- run.go — Sets up namespaces, forks child, waits
- rootfs.go — Mounts overlayfs, calls pivot_root
- proc.go — Mounts /proc /sys /dev inside container
- cleanup.go — Unmounts, removes cgroup, cleans scratch dirs

### cgroups/

Writes resource limits to cgroup filesystem.
Detects v1 vs v2 at runtime.

- cgroup.go — Apply(pid, limits), Destroy(id)

### network/

Creates veth pair, configures IPs, NAT rules.
Currently shells out to `ip` commands (netlink later).

- network.go — Setup(containerPID, config), Teardown(id)

### image/

Manages rootfs tarballs. Simple directory per image.
No registry pull yet — manual import only.

- image.go — Unpack(tarPath, dest), Exists(name)

---

## Execution Flow

Mentioned in Docs/Diagrams/ExecutionFlow.md

---

## Key Design Decisions

**Re-exec pattern**: Go runtime initializes before you can set namespaces,
so the binary re-execs itself with a sentinel env var to run namespace-setup
code before Go runtime takes over. Same trick runc uses.

**No daemon**: No background daemon. Each run is a self-contained process tree.

**Overlayfs over plain copy**: Lets multiple containers share one base image
without duplicating filesystem. Teaches the actual layer model.

**Cgroup v2 first**: Unified hierarchy. v1 is legacy. Most modern kernels
default v2.

---

## What's Intentionally Missing

- Image registry pull (OCI spec is a rabbit hole — add later)
- User namespaces (adds complexity; learn namespaces first)
- seccomp / capabilities filtering
- Port forwarding
- Container-to-container networking
