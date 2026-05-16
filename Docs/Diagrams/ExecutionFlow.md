# Execution Flow

This document explains how a container is created and executed inside
Mini Docker.

The most important thing to understand:

> A container is NOT a virtual machine.
>
> A container is just a Linux process running inside isolated kernel views.

The host kernel is always shared.

---

# HOST CONTEXT

```text
USER
 │
 │ runs:
 │ minidocker run alpine /bin/sh
 ▼

┌─────────────────────────────────────────────┐
│ HOST PROCESS                                │
│ minidocker (parent process)                 │
│                                             │
│ Responsibilities:                           │
│ - parse CLI args                            │
│ - prepare root filesystem                   │
│ - spawn isolated child process              │
│ - apply cgroups                             │
│ - setup networking                          │
│ - cleanup resources                         │
└─────────────────────────────────────────────┘
                    │
                    │ image.Prepare()
                    │ rootfs.Mount()
                    │
                    ▼

═══════════════════════════════════════════════════════
 OVERLAYFS SETUP (still on host)
═══════════════════════════════════════════════════════

lowerdir = alpine image (read-only)
upperdir = writable container layer
merged   = final container filesystem

The merged directory becomes the future container root filesystem.

At this point:
- process is still on the host
- host namespaces are still active
- no container exists yet

                    │
                    │ exec.Command("/proc/self/exe", "init")
                    │ + Cloneflags
                    ▼

═══════════════════════════════════════════════════════
 NEW PROCESS CREATED WITH NEW NAMESPACES
═══════════════════════════════════════════════════════

Host process tree:

shell
 └── minidocker parent
      └── minidocker child

The child process is cloned into new namespaces:

- PID namespace
- Mount namespace
- UTS namespace
- Network namespace
- IPC namespace

Important:
- the child is still visible on the host
- namespaces change what the child can SEE
- namespaces do NOT create a VM

                    │
                    │ child process starts
                    │
                    ├── cgroups.Apply(childPID)
                    │
                    ├── network.Setup(childPID)
                    │
                    ▼

═══════════════════════════════════════════════════════
 RE-EXEC PATTERN
═══════════════════════════════════════════════════════

The child re-execs the SAME binary:

exec.Command("/proc/self/exe", "init")

Why this exists:

- Go runtime creates threads very early
- namespace operations require controlled process state
- container runtimes solve this using re-exec

This is similar to how runc works internally.

                    │
                    ▼

# CONTAINER CONTEXT

┌─────────────────────────────────────────────┐
│ CONTAINER INIT PROCESS                      │
│ (already inside namespaces)                 │
│                                             │
│ code path:                                  │
│                                             │
│ if arg == "init" {                          │
│     containerInit()                         │
│ }                                           │
└─────────────────────────────────────────────┘
                    │
                    │ pivot_root()
                    ▼

═══════════════════════════════════════════════════════
 FILESYSTEM SWITCH
═══════════════════════════════════════════════════════

Before:

    /  -> host filesystem

After:

    /  -> container merged overlayfs

The process now sees the overlay filesystem as `/`
inside its mount namespace.

Important:
- this is NOT a separate kernel
- this is filesystem isolation
- mount namespace + pivot_root work together

                    │
                    │ mount /proc /sys /dev
                    ▼

═══════════════════════════════════════════════════════
 CONTAINER ENVIRONMENT SETUP
═══════════════════════════════════════════════════════

Fresh virtual filesystems are mounted:

- /proc
- /sys
- /dev

Why this matters:

- `/proc` must be remounted for correct PID visibility
- tools like `ps` depend on this
- mounts are isolated by the mount namespace

Hostname is also isolated:

syscall.Sethostname(...)

                    │
                    ▼

═══════════════════════════════════════════════════════
 PROCESS REPLACEMENT
═══════════════════════════════════════════════════════

Final step:

syscall.Exec("/bin/sh", ...)

IMPORTANT:

exec() does NOT create a new process.

It REPLACES the current process image.

SAME PROCESS
SAME PID
NEW EXECUTABLE

Before:

    PID 1 -> minidocker init

After:

    PID 1 -> /bin/sh

PID 1 inside containers is special:
- receives signals differently
- responsible for zombie reaping

                    │
                    ▼

═══════════════════════════════════════════════════════
 FINAL STATE
═══════════════════════════════════════════════════════

Host sees:

shell
 └── minidocker parent
      └── /bin/sh

Container sees:

PID 1 -> /bin/sh

Different views exist because of PID namespaces.

---

# FINAL MENTAL MODEL

Mini Docker itself is NOT the container.

Mini Docker is:

- a process orchestrator
- a namespace manager
- a filesystem switcher
- a resource configurator

The actual container is simply:

    a Linux process running inside isolated kernel views

Isolation comes from:
- namespaces
- cgroups
- mount isolation
- network isolation

NOT from virtualization.
```
