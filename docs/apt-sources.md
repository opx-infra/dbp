# Apt Sources

*dbp* supports OpenSwitch's own Apt sources, along with any custom Apt sources required. *dbp* **always** gives priority to locally built packages, even if the version is lower.

## OpenSwitch Apt sources

*dbp* loads [OpenSwitch apt sources](http://deb.openswitch.net) into `/etc/apt/sources.list.d/20extra.list`if no other sources are specified.

```bash
$ dbp -v build
[INFO] Loaded extra sources:
deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx
```

Using the `--dist` flag modifies `stretch` in the above list, while using the `--release` flag modifies `unstable`.

```bash
$ dbp --dist bionic --release stable -v build
[INFO] Loaded extra sources:
deb     http://deb.openswitch.net/bionic stable opx opx-non-free
deb-src http://deb.openswitch.net/bionic stable opx
```

Using the `--no-extra-sources` flag removes these sources entirely.

## Custom Apt sources

There are many ways to specify custom Apt sources. In order of precedence:

* The `--extra-sources` flag
* The `$EXTRA_SOURCES` environment variable
* The `./extra_sources.list` file
* The `$HOME/.extra_sources.list` file

For example, fill `~/.extra_sources.list` with

```
deb http://deb.notarealsite.net/ sid main
```

and *dbp* will use it.

```bash
$ dbp -v build
[INFO] Loaded extra sources:
deb http://deb.notarealsite.net/ sid main
```
