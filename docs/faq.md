# Frequently Asked Questions

## What is happening with opx-build?

OpenSwitch will continue to support the [`opx_build`](https://github.com/open-switch/opx-build/blob/master/scripts/opx_build) script from the [`opx-build`](https://github.com/open-switch/opx-build) repository for building OPX Debian packages and the OPX installer. `opx-build` currently supports building Debian packages for both Debian Jessie and Debian Stretch.

!!! tip
    opx-build will be continue to be supported as long as OPX ships on Debian Stretch.

## How do I build an unstripped and unoptimized binary?

Set `DEB_BUILD_OPTIONS='nostrip noopt debug'`. You can use it with `build`, `gbp buildpackage`, and even `fakeroot debian/rules binary`.

```bash
$ dbp shell
build@stretch:/mnt$ DEB_BUILD_OPTIONS='nostrip noopt debug' build opx-logging
build@stretch:/mnt$ ll pool/stretch-amd64/opx-logging/python-opx-logging*
-rw-r--r-- 1 build dialout  12K Aug 29 22:03 pool/stretch-amd64/opx-logging/python-opx-logging_2.1.1_amd64.deb

build@stretch:/mnt/opx-logging$ cd opx-logging/
build@stretch:/mnt/opx-logging$ DEB_BUILD_OPTIONS='nostrip noopt debug' gbp buildpackage
build@stretch:/mnt/opx-logging$ ll ../python-opx-logging*
-rw-r--r-- 1 build dialout 12K Aug 29 22:04 ../python-opx-logging_2.1.1_amd64.deb

# Normally you'd get this
build@stretch:/mnt/opx-logging$ cd ..
build@stretch:/mnt$ rm -rf pool/
build@stretch:/mnt$ build opx-logging
build@stretch:/mnt$ ll pool/stretch-amd64/opx-logging/python-opx-logging*
-rw-r--r-- 1 build dialout  12K Aug 29 22:07 pool/stretch-amd64/opx-logging/python-opx-logging-dbgsym_2.1.1_amd64.deb
-rw-r--r-- 1 build dialout 4.2K Aug 29 22:07 pool/stretch-amd64/opx-logging/python-opx-logging_2.1.1_amd64.deb
```

## How do I generate a debian/changelog entry?

```bash
$ dbp shell -c "cd src/; gbp dch -R"
```

An editor will be opened where you can edit the changelog entry. Simply commit the changes and raise a pull request!

## How do I build in parallel?

You can override the default parallel level with a build option. Set `N` to any number >0.

```bash
DEB_BUILD_OPTIONS='parallel=N'
```

## What is the `build` program?

`build` is a simple shell script which `cd`s into the source directory and runs `gbp buildpackage` with the [`--git-export-dir`](http://honk.sigxcpu.org/projects/git-buildpackage/manual-html/man.gbp.buildpackage.html) flag. This sorts build artifacts into per-source package directories.

```bash
#!/usr/bin/env bash
set -euo pipefail

if [[ $# -gt 0 ]]; then
  cd "$1"
  NAME="$1"
  shift
else
  NAME="$(basename "$(pwd)")"
fi

gbp buildpackage --git-export-dir="/mnt/pool/${DIST}-${ARCH}/${NAME}" "$@"
```

## What is the `install-build-deps` program?

`install-build-deps` is a shell script which indexes local packages, adds your extra sources, then installs build dependencies. It is run automatically when using `gbp buildpackage` or `build`. Skip running it with `dbp build --gbp="--git-prebuild=':'"`.

## Why am I receiving `unable to find user build`?

Docker could be too slow (or it could be storage). Either way, the container is taking too long to start. Try running `dbp run` first before running a build or shell.
