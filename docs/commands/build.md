# Commands: *dbp build*

Out of the box, *dbp build* reads all directories and builds a graph of Debian build dependencies, which it traverses while building each Debian package. These packages are sorted into a pool and used for subsequent builds and ultimately publishing.

## Build a single package in a non-persistent container

```bash
$ dbp build src/
```

* Builds artifacts for the default Debian distribution
* Uses packages found in `./pool/stretch-amd64` as build dependencies
* Deposits artifacts in `./pool/stretch-amd64/src/`
* If workspace container does not exist, a container is created for this build and destroyed after
* If the workspace container already exists, it is used for the build and *not* destroyed after

## Build multiple repositories

With no directories specified, dbp will build in build dependency order.

```bash
$ dbp build
```

Otherwise, directories can be manually specified.

```bash
$ dbp build src new-src amazing-src
```

## Building unstripped and unoptimized binaries

Use the `--debug` flag with *dbp build* or *dbp shell*. If already in a shell session, just run `export DEB_BUILD_OPTIONS='nostrip noopt debug'`. It works with `build`, `gbp buildpackage`, and even `fakeroot debian/rules binary`.

Without `--debug`:
```bash hl_lines="3 4"
build@stretch:/mnt$ ll pool/stretch-amd64/opx-logging | awk '{print $5, $9}' | column -t
13K   libopx-logging-dev_2.1.1_amd64.deb
85K   libopx-logging1-dbgsym_2.1.1_amd64.deb
10K   libopx-logging1_2.1.1_amd64.deb
905   opx-logging_2.1.1.dsc
12K   opx-logging_2.1.1.tar.gz
28K   opx-logging_2.1.1_amd64.build
6.9K  opx-logging_2.1.1_amd64.buildinfo
3.1K  opx-logging_2.1.1_amd64.changes
3.6K  opx-logging_2.1.1_amd64.deb
12K   python-opx-logging-dbgsym_2.1.1_amd64.deb
4.2K  python-opx-logging_2.1.1_amd64.deb
```

With `--debug`:
```bash hl_lines="9"
$ dbp --debug shell

build@stretch:/mnt$ echo "Options: $DEB_BUILD_OPTIONS"
Options: nostrip noopt debug

build@stretch:/mnt$ build opx-logging
build@stretch:/mnt$ ll pool/stretch-amd64/opx-logging | awk '{print $5, $9}' | column -t
93K   libopx-logging-dev_2.1.1_amd64.deb
85K   libopx-logging1_2.1.1_amd64.deb
905   opx-logging_2.1.1.dsc
12K   opx-logging_2.1.1.tar.gz
28K   opx-logging_2.1.1_amd64.build
6.4K  opx-logging_2.1.1_amd64.buildinfo
2.6K  opx-logging_2.1.1_amd64.changes
3.6K  opx-logging_2.1.1_amd64.deb
12K   python-opx-logging_2.1.1_amd64.deb
```

## Build against an OPX release

```bash
$ dbp --release 3.0.0 build
```

## Build against a Debian distribution

```bash
$ dbp --dist bionic build
```

## Build against extra apt sources

*dbp* will read from the following list of inputs for extra apt sources. These sources must be in standard sources.list format.

1. `--extra-sources` argument
1. `EXTRA_SOURCES` environment variable
1. `./.extra_sources.list` file
1. `~/.extra_sources.list` file

For example, fill `~/.extra_sources.list` with

```
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx
```

and *dbp build* will search OpenSwitch for build dependencies.

```bash
$ dbp -v build src
INFO:dbp:Loaded extra sources:
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx
```

If no sources are found, the default OPX sources are used.

## Exclude custom apt sources

Also excludes default OPX sources. Useful for ensuring a complete local build is possible.

```bash
$ dbp --no-extra-sources build
```

## Pass additional `git-buildpackage` options

For example, skip building when tagging by passing the correct flag.

```bash
dbp build src --gbp="--git-tag-only"
```

## Build in (more) parallel

You can override the default parallel level with a build option. Set `N` to any number >0.

```bash
DEB_BUILD_OPTIONS='parallel=N'
```
