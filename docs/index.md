# Docker-buildpackage Documentation

Docker-buildpackage is a tool for building Debian packages.

## Prerequisites

1. Python 3.5+
1. Docker

## Install

```bash
pip3 install dbp
```

## Usage

By default, `dbp` uses OpenSwitch apt sources[^1].

* `dbp build` runs an out-of-tree build and stores build artifacts in `./pool/` for easy publishing
* `dbp shell` launches an interactive bash shell in the development environment container
* `dbp run` starts a persistent container in the background
* `dbp rm` removes the persistent container from the background

Both `dbp build` and `dbp shell` use temporary containers if no container exists.

[^1]:
    ```
    deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
    deb-src http://deb.openswitch.net/stretch unstable opx
    ```
