# Frequently Asked Questions

## What is happening with opx-build?

OpenSwitch will continue to support the [*opx_build*](https://github.com/open-switch/opx-build/blob/master/scripts/opx_build) script from the [*opx-build*](https://github.com/open-switch/opx-build) repository for building OPX Debian packages and the OPX installer. *opx-build* currently supports building Debian packages for both Debian Jessie and Debian Stretch.

!!! tip
    opx-build will be continue to be supported as long as OPX ships on Debian Stretch.

## How do I generate a debian/changelog entry?

```bash
$ dbp shell -c "cd src/; gbp dch -R"
```

An editor will be opened where you can edit the changelog entry. Simply commit the changes and raise a pull request!

## What is the `install-build-deps` program?

`install-build-deps` is a shell script which indexes local packages, adds your extra sources, then installs build dependencies. It is run automatically when using `gbp buildpackage` or `build`. Skip running it with `dbp build --gbp="--git-prebuild=':'"`.

## What is the `pool-packages` program?

`pool-packages` reads a Debian changes file and moves all associated files into `./pool/${DIST}-${ARCH}/src`. It is useful for sharing packages.

## Why am I receiving `unable to find user build`?

Docker could be too slow (or it could be storage). Either way, the container is taking too long to start. Try running *dbp run* first before running a build or shell.
