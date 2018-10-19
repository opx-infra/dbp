# Building ONIE Installers

[*opx-build*](https://github.com/open-switch/opx-build/tree/installer/3.0.0) is used for building ONIE installers. It is pre-installed in the development environment.

```bash
# Clone installer blueprint
$ git clone https://github.com/open-switch/opx-onie-installer

# Run installer builder script in container
$ dbp shell -c 'opx_rel_pkgasm.py \
  -b opx-onie-installer/release_bp/OPX_dell_base_stretch.xml'
```

## Usage

```bash
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
