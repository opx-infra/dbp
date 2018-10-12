# Docker-buildpackage Documentation

Docker-buildpackage is a tool for building Debian packages.

!!! note
    See why *dbp* is replacing *opx-build* [here](comparison.md).

OpenSwitch's development workflow is outlined [here](development-workflow.md).

## Install

```bash
$ pip3 install --upgrade dbp
```

For more information, see [installation](installation.md).

## Support

For bug reports and feature requests, file an issue on [GitHub](https://github.com/opx-infra/dbp/issues/new/choose). For other support, drop by the [*dbp* support channel](https://chat.openswitch.net/channel/dbp) and we'll help you out.

## Usage

1. Clone any number of OpenSwitch repositories
2. Build!

```bash
$ git clone https://github.com/open-switch/opx-logging
$ dbp build
```

### Commands

* [*dbp build*](commands/build.md) runs an out-of-tree build and stores build artifacts in `./pool/` for easy publishing
* [*dbp shell*](commands/shell.md) launches an interactive bash shell in the development environment container
* [*dbp run*](commands/run.md) starts a persistent container in the background
* [*dbp rm*](commands/rm.md) removes the persistent container from the background

!!! tip
    Read up on *dbp* and Apt sources [here](apt-sources.md).

## Building OpenSwitch Installers

*dbp* can be used to build an OpenSwitch installer.

```bash
$ git clone https://github.com/open-switch/opx-onie-installer
$ dbp shell -c 'opx_rel_pkgasm.py \
  -b opx-onie-installer/release_bp/OPX_dell_base_stretch.xml'
```

Usage:

```bash
$ dbp shell -c 'opx_rel_pkgasm.py --help'
usage: opx_rel_pkgasm.py [-h] -b B [-n N] [-s S] [-v V]
                         [--build-info BUILD_INFO] [--build-url BUILD_URL]
                         [--vcs-url VCS_URL] [--vcs-revision VCS_REVISION]
                         [-d DIST]

optional arguments:
  -h, --help            show this help message and exit
  -b B                  specify location of release blue-print
  -n N                  specify build number of release
  -s S                  specify release number suffix
  -v V                  specify verbosity level
  --build-info BUILD_INFO
                        specify location of build-info json output
  --build-url BUILD_URL
  --vcs-url VCS_URL
  --vcs-revision VCS_REVISION
  -d DIST, --dist DIST  Distribution to build
```
