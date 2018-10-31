# Building OpenSwitch packages
Docker-buildpackage supports both interactive and non-interactive package builds.

For this page, assume the following workspace is present.
```bash
$ git clone https://github.com/open-switch/opx-logging
$ git clone https://github.com/open-switch/opx-common-utils
```

## Non-interactive package builds
Build specific repositories by naming them.
```bash
$ dbp build opx-logging
```

Packages will be deposited in your workspace after a successful build. Use `pool-packages` to pool packages into per-repository directories for easy sharing.
```bash
$ ls *.changes
opx-logging_2.1.1_amd64.changes
$ dbp shell -c 'pool-packages *.changes'
$ ls pool/stretch-amd64/opx-logging/*.changes
pool/stretch-amd64/opx-logging/opx-logging_2.1.1_amd64.changes
```

Use *dbp build* without any arguments to build a full workspace. This runs `cd $repo; gbp buildpackage` for each repository with a `debian/control` file.
```bash
$ dbp build
```

The `--print` flag can be used to preview the build order.
```bash
$ dbp build --print
opx-logging opx-common-utils
```

## Interactive package builds
When developing a package, use *dbp run* to launch a persistent development container in the background. It can then be entered with *dbp shell* and removed with *dbp rm*.
```bash
$ dbp run
$ dbp shell
$ dbp rm
```

!!! tip
    Running `dbp build opx-logging` is functionally equivalent to
    ```bash
    $ dbp shell -c 'cd opx-logging; gbp buildpackage'
    ```

Running *dbp shell* launches a Bash session, where the standard Debian workflow is available.

```bash
$ dbp shell
build@stretch:/mnt$ cd opx-logging

# Install build dependencies and build the package
build@stretch:/mnt/opx-logging$ gbp buildpackage

# On failed builds, avoid the long gbp build time by quickly rebuilding
build@stretch:/mnt/opx-logging$ fakeroot debian/rules binary

# Manually clean up
build@stretch:/mnt/opx-logging$ fakeroot debian/rules clean
build@stretch:/mnt/opx-logging$ exit
```

*dbp run* creates a container based on the current user, directory, and shell process id. To always use the same container name (or start a shell in someone else's), use the `$CNAME` environment variable.

```bash
$ dbp -vv rm
[DEBUG] Running docker inspect user-dbp-demo-22222

$ CNAME=CONSTANT dbp -vv rm
[DEBUG] Running docker inspect CONSTANT
```

## Debug builds
Build unstripped, unoptimized packages with the `--debug` flag. Both *dbp shell* and *dbp build* support it.
```bash
$ dbp --debug build opx-logging
$ dbp --debug shell -c 'cd opx-logging; gbp buildpackage'
```

## Pass additional Git-buildpackage options
Use the `--gbp` flag with *dbp build*.
```bash
dbp build opx-logging --gbp="--git-tag-only"
```
