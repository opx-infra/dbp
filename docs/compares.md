# Docker-buildpackage and *opx-build*

*dbp* maintains some backwards compatibility with *opx-build*. `opx_build` does not work without *pbuilder*, and `opx_run` does not work without Docker-in-Docker. The [*installer/3.0.0*](https://github.com/open-switch/opx-build/tree/installer/3.0.0) tag of *opx-build* is installed into `/usr/local/opx-build` and `/usr/local/opx-build/scripts` has been added to the `$PATH`.

```bash
$ dbp shell -c 'which opx_rel_pkgasm.py'
/usr/local/opx-build/scripts/opx_rel_pkgasm.py
```

## Commands

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

$ dbp shell -c 'opx-build/scripts/opx_rel_pkgasm.py --dist stable \
  -b opx-onie-installer/release_bp/OPX_dell_base_stretch.xml'
```

## Speed

This should definitely be better tested.

```bash
$ time opx-build/scripts/opx_run opx-build opx-logging
opx-build/scripts/opx_run opx_build opx-logging  0.05s user 0.03s system 0% cpu 3:21.92 total

$ time dbp build opx-logging
dbp build opx-logging  0.50s user 0.16s system 1% cpu 49.015 total
```
