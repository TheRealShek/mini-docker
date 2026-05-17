# Development

## Local Setup

To continue development, prepare a small Alpine root filesystem on a Linux host:

```bash
mkdir -p /tmp/alpine-rootfs
curl -L https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz \
  | sudo tar -xz -C /tmp/alpine-rootfs
```

Verify that the root filesystem was extracted:

```bash
ls /tmp/alpine-rootfs
```
