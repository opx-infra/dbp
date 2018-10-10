# Commands: *dbp rm*

*dbp rm* removes and container started by [*dbp run*](run.md) from the same directory. If no container exists, nothing is done.

```bash
$ dbp -vv rm
[INFO] Loaded extra sources:
deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx

[DEBUG] Running docker inspect theucke-dbp-dbp
[DEBUG] Running docker rm -f theucke-dbp-dbp
```
