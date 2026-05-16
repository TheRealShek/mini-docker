# Prerequisites — What You Should Know Before Writing a Single Line

This isn't a checklist to intimidate you. It's stuff I wish I'd actually
understood before jumping in, because if you don't, you'll spend more time
confused about _why something isn't working_ than actually learning containers.

You don't need to master all of this. But you need to be comfortable enough
that when you see a syscall or a mount flag in someone's code, you're not
completely lost.

---

> You can Manually go to understand Topic by Topic or Read this-
> [What Your Container Runtime Actually Does When You Type docker run](https://therealshek.medium.com/what-your-container-runtime-actually-does-when-you-type-docker-run-089f13fbdccb)

## 1. How Linux Creates Processes

Before containers make sense, processes need to make sense.

When you run a program, Linux does two things — `fork` creates a copy of the
current process, and `exec` replaces that copy with the new program. That's it.
Every process on your system came from this. Your shell, your browser, everything.

What you should understand:

- Every process has a parent. If the parent dies before the child, the child
  gets reparented to PID 1 (init/systemd).
- Zombie processes — a child that's exited but whose parent hasn't called
  `wait()` yet. You'll create zombies by accident in this project.
- What `/proc/<pid>/` actually contains. Spend 10 minutes just poking around
  `/proc/self/` in your terminal. It'll click.

Where to look: `man 2 fork`, `man 2 execve`, `man 5 proc`

---

## 2. Linux Namespaces — The Whole Point

Namespaces are the actual mechanism behind container isolation. Everything else
is built on top of them. You need to understand what each one does before you
touch code, because otherwise you're just copy-pasting flags without knowing why.

The ones that matter for this project:

**PID namespace** — Processes inside see their own PID tree. The first process
inside becomes PID 1. From the host, you can still see them with their real PIDs.

**Mount namespace** — Filesystem mounts are isolated. You can mount and unmount
things inside without affecting the host. This is what makes `pivot_root` useful.

**UTS namespace** — Just the hostname and domain name. Lets the container have
its own hostname without touching the host's.

**Network namespace** — Separate network stack. Own interfaces, routing table,
iptables rules. This is why a container can have its own IP.

**IPC namespace** — Isolates System V IPC and POSIX message queues. Less
exciting, but needed for full isolation.

You don't need to memorize flags. You need to understand _what gets isolated_
in each namespace and _why that matters_. Read:
`man 7 namespaces` — seriously, it's actually readable.

---

## 3. Syscalls from Go

Go doesn't hide syscalls from you. The `syscall` package and
`golang.org/x/sys/unix` expose them directly. You'll use them constantly.

What you need to be comfortable with:

- Reading a syscall signature and mapping it to what Go exposes
- `SysProcAttr` on `exec.Cmd` — this is how you set namespace flags, set the
  child's controlling terminal, kill the child if the parent dies, etc.
- The difference between `syscall.Exec` (replaces current process) and
  starting a subprocess with `exec.Cmd`

You don't need to be a syscall expert. But when you see:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
}
```

You should know exactly what that's doing and why.

---

## 4. Filesystems and Mounting

This one trips people up more than anything else.

Understand these concepts before Phase 2:

**What a mount actually is** — It's not copying files. It's attaching a
filesystem (or directory) to a point in the directory tree. The original
directory is hidden behind it, not deleted.

**Bind mounts** — Mounting a directory onto another directory. Same files,
different path. You'll use this constantly. `mount --bind /some/dir /other/dir`.

**Mount propagation** — Whether mounts inside a namespace are visible outside.
`MS_PRIVATE`, `MS_SHARED`, `MS_SLAVE`. You'll need this so your container's
mounts don't bleed into the host.

**pivot_root vs chroot** — `chroot` changes what `/` looks like but doesn't
change the actual root mount. Someone with enough access can escape it.
`pivot_root` actually swaps the root filesystem at the kernel level. Use
`pivot_root`. Read `man 2 pivot_root` before you need it, not during.

**Why `/proc` has to be mounted manually** — When you create a new PID
namespace, `/proc` still shows the host's processes until you mount a fresh
one. This is confusing the first time you see it.

---

## 5. Cgroups (Control Groups)

Cgroups limit what resources a process (and its children) can use — CPU,
memory, I/O, etc. They're a separate mechanism from namespaces. Namespaces
control _what you can see_, cgroups control _what you can use_.

The thing to know upfront: there are two versions, v1 and v2, and they work
differently. Most modern Linux systems (Fedora included) use v2 by default.

**cgroup v1** — Each resource type has its own hierarchy under
`/sys/fs/cgroup/memory/`, `/sys/fs/cgroup/cpu/`, etc.

**cgroup v2** — Unified hierarchy. Everything lives under `/sys/fs/cgroup/`.
You write to files like `memory.max`, `cpu.max`, `cgroup.procs`.

Check which one you have:

```
ls /sys/fs/cgroup/
```

If you see `cgroup.controllers`, you're on v2. Learn v2. Don't waste time on
v1 tutorials that show paths that don't exist on your machine.

---

## 6. Basic Networking Concepts

You don't need this until Phase 5, but it helps to have a rough model in your
head from the start.

**Network namespace** — Completely separate network stack. No interfaces by
default except loopback.

**veth pair** — Virtual ethernet cable. Two ends — put one in the container's
network namespace, keep one on the host. Traffic in one end comes out the other.

**Bridge interface** — A virtual switch on the host. Connect multiple veth host
ends to it, and containers can talk to each other through it.

**NAT / iptables** — How container traffic reaches the internet. The host
masquerades the container's private IP with its own public IP.

You don't need to configure any of this yet. Just know the shape of it.

---

## 7. Go Fundamentals You'll Actually Use

Assuming you know Go basics. The parts that come up constantly in this project:

- `os/exec` — spawning child processes, inheriting file descriptors
- `syscall` and `golang.org/x/sys/unix` — raw Linux interface
- `os.MkdirAll`, `os.RemoveAll` — filesystem setup and teardown
- Error wrapping — you'll be debugging deep syscall errors; `fmt.Errorf("%w", err)`
  and knowing how to unwrap matters
- `defer` for cleanup — you'll have a lot of mounts and cgroups to clean up on exit

The re-exec pattern is the one Go-specific thing that's genuinely non-obvious.
Because Go's runtime starts goroutines before `main()` runs, you can't safely
call `unshare()` or set certain namespace flags after the runtime has
initialized. The fix is to have the binary re-execute itself with a special
environment variable, and in that second execution, do the namespace work
before Go's runtime fully starts. You'll implement this in Phase 1. Just
know it's coming.

---

## How to Verify You're Ready

You should be able to answer these without Googling:

- What happens to a child process if the parent exits?
- What does `CLONE_NEWPID` actually isolate?
- What's the difference between a bind mount and a regular mount?
- Why can't you just use `chroot` for filesystem isolation?
- Where on your filesystem do cgroup controls live?

If two or three of these feel shaky, go read the relevant section again before
starting. An hour of reading now saves three hours of debugging later.
