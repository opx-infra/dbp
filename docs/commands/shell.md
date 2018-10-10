# Commands: *dbp shell*

*dbp shell* launches an interactive shell in a Debian development environment. This environment is provided by [*opx-infra/gbp-docker*](https://github.com/opx-infra/gbp-docker). If a container were previously launched with [*dbp run*](run.md) in the same directory, *dbp shell* would launch the shell in the existing container. This existing container is not removed when the shell exits. If no container exists, a new container is created and then destroyed after.

```bash
$ dbp -vv shell
[INFO] Loaded extra sources:
deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx

[DEBUG] Running docker inspect theucke-dbp-dbp
[DEBUG] Running docker inspect theucke-dbp-dbp
[DEBUG] Running docker run -d -it --name=theucke-dbp-dbp --hostname=stretch -v=/neteng/theucke/opx/dbp:/mnt -v=/home/theucke/.gitconfig:/etc/skel/.gitconfig:ro -e=UID=9438 -e=GID=3000 -e=TZ=US/Pacific-New -e=DEBFULLNAME=Dell EMC -e=DEBEMAIL=ops-dev@lists.openswitch.net -e=EXTRA_SOURCES=deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx
 opxhub/gbp:v1.0.5-stretch-dev bash -l
[DEBUG] Running docker exec -it --user=build -e=UID=9438 -e=GID=3000 -e=TZ=US/Pacific-New -e=DEBFULLNAME=Dell EMC -e=DEBEMAIL=ops-dev@lists.openswitch.net -e=EXTRA_SOURCES=deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx
 theucke-dbp-dbp bash -l
build@stretch:/mnt$ exit
logout
[DEBUG] Running docker inspect theucke-dbp-dbp
[DEBUG] Running docker rm -f theucke-dbp-dbp
```

## Develop inside a persistent development container

Using [*dbp run*](run.md) launches a persistent development container. This container will only be explicitly removed when [*dbp rm*](rm.md) is run in the same directory. You can use *dbp shell* to enter this container, or `dbp shell -c 'cmd'` to run a command non-interactively.

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
