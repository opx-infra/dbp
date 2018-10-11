# Docker-buildpackage and *opx-build*

*dbp* is replacing *opx-build* to support more OS distributions, increase developer productivity, and lower operational overhead.

*opx-build* contains scripts for running the OPX Docker image [*opxhub/build*](https://hub.docker.com/r/opxhub/build) (`opx_run`), building Debian packages with [*Git-buildpackage*](https://honk.sigxcpu.org/piki/projects/git-buildpackage/) (`opx_build`), and [building an OpenSwitch installer](https://github.com/open-switch/opx-build/blob/master/scripts/opx_rel_pkgasm.py) (`opx_rel_pkgasm.py`).

#### On *opx-build* and build isolation

*opx-build* uses [*git-buildpackage*](https://honk.sigxcpu.org/piki/projects/git-buildpackage/)/[*git-pbuilder*](https://wiki.debian.org/git-pbuilder)/[*cowbuilder*](https://wiki.debian.org/cowbuilder) for building packages. This provides clean environments for every build, but is also rather redundant with the Docker image we provide. *dbp* resolves this issue by removing *git-pbuilder*/*cowbuilder* and running builds in the container directly. This also brings the added benefit of simplified environment management, as we can build one image per distribution, instead of one image containing every distribution.

#### On *opx-build* and container persistence

With *opx-build*, the container is removed after every build. Developers are forced to manually run containers and execute shell sessions within them. *dbp* is built from the ground up around this workflow, supporting persistent containers and multiple shell sessions.

## *dbp* and backwards compatibility

*dbp* maintains *some* backwards compatibility with *opx-build*. `opx_build` does not work without *pbuilder*, and `opx_run` does not work without Docker-in-Docker. The [*installer/3.0.0*](https://github.com/open-switch/opx-build/tree/installer/3.0.0) tag of *opx-build* is installed into `/usr/local/opx-build` and `/usr/local/opx-build/scripts` has been added to the `$PATH` to support installer building.

```bash
$ dbp shell -c 'which opx_rel_pkgasm.py'
/usr/local/opx-build/scripts/opx_rel_pkgasm.py
```

## Command Comparison

**Build a single package**

```bash
$ opx-build/scripts/opx_run opx_build src1

$ dbp build src1
```

**Build all packages in a directory**

```bash
$ opx-build/scripts/opx_run opx_build src1 src2 src3...

$ dbp build
```

**Build an OPX installer**

```bash
$ opx-build/scripts/opx_run opx_rel_pkgasm.py --dist stable \
  -b opx-onie-installer/release_bp/OPX_dell_base_stretch.xml

$ dbp shell -c 'opx_rel_pkgasm.py --dist stable \
  -b opx-onie-installer/release_bp/OPX_dell_base_stretch.xml'
```
