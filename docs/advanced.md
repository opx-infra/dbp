# Advanced Usage

Here are some fun things you can do with `dbp`.

## Build a single package in a non-persistent container

```bash
dbp build src/
```

* Builds artifacts for the default Debian distribution
* Uses packages found in `./pool/stretch-amd64` as build dependencies
* Deposits artifacts in `./pool/stretch-amd64/src/`
* If workspace container does not exist, a container is created for this build and destroyed after
* If the workspace container already exists, it is used for the build and *not* destroyed after

## Build multiple repositories

With no directories specified, dbp will build in build dependency order.

```bash
dbp build
```

Otherwise, directories can be manually specified.

```bash
dbp build src new-src amazing-src
```

## Build against an OPX release

```bash
dbp --release 3.0.0 build
```

## Build against a Debian distribution

```bash
dbp --dist bionic build
```

## Build against extra apt sources

`dbp` will read from the following list of inputs for extra apt sources. These sources must be in standard sources.list format.

1. `--extra-sources` argument
1. `EXTRA_SOURCES` environment variable
1. `./.extra_sources.list` file
1. `~/.extra_sources.list` file

For example, fill `~/.extra_sources.list` with

```
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx
```

and `dbp build` will search OpenSwitch for build dependencies.

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
dbp --no-extra-sources build
```

## Develop inside a persistent development container

Using `dbp run` launches a persistent development container. This container will only be explicitly removed when `dbp rm` is run in the same directory. You can use `dbp shell` to enter this container, or `dbp shell -c 'cmd'` to run a command non-interactively.

```bash
dbp run
dbp shell

# Now we are inside the container (denoted by $ prompt)
$ cd src/

# Install build dependencies and build the package
$ gbp buildpackage

# Only install build dependencies
$ install-build-deps

# On failed builds, avoid the long gbp build time by quickly rebuilding
$ fakeroot debian/rules build

# Manually clean up
$ fakeroot debian/rules clean

# Add a new source for build dependencies by appending to the env var
$ export EXTRA_SOURCES="$EXTRA_SOURCES
deb http://deb.openswitch.net/stretch 3.0.0 opx opx-non-free"

# Run gbp buildpackage again to do a clean build, but this time skip installing build deps
$ gbp buildpackage --git-prebuild=':'

# Run an "official" build (what the CI runs)
$ cd /mnt
$ build src

# Build an unstripped, unoptimized binary (this also works with gbp and debian/rules)
$ DEB_BUILD_OPTIONS='nostrip noopt debug' build src

# Exit the container
$ exit

# Remove the container when finished (or use `dbp shell/build` again to re-enter the same container)
dbp rm
```

!!! note
    Packages are only indexed for build dependencies in `pool/${DIST}-${ARCH}/`. Building with `build ./src` automatically deposits build artifacts into the correct directory for indexing. You can also simply copy `.deb` files into the directory.

## Pass additional `git-buildpackage` options

For example, skip building when tagging by passing the correct flag.

```bash
dbp build src --gbp="--git-tag-only"
```

## Pull any Docker image updates

```bash
dbp pull
dbp -d buster pull
```
