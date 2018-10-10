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

* [*dbp build*](commands/build.md) runs an out-of-tree build and stores build artifacts in `./pool/` for easy publishing
* [*dbp shell*](commands/shell.md) launches an interactive bash shell in the development environment container
* [*dbp run*](commands/run.md) starts a persistent container in the background
* [*dbp rm*](commands/rm.md) removes the persistent container from the background

*dbp* uses [OpenSwitch apt sources](http://deb.openswitch.net)[^1] if no other sources are specified.
Both *dbp build* and *dbp shell* use temporary containers if no container exists.

[^1]:
    ```
    deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
    deb-src http://deb.openswitch.net/stretch unstable opx
    ```
