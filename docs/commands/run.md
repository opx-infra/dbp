# Commands: *dbp run*

*dbp run* creates and starts a development container that can be used by both [*dbp build*](build.md) and [*dbp shell*](shell.md). It must be manually removed with [*dbp rm*](rm.md).

```bash
$ dbp -vv run
[INFO] Loaded extra sources:
deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx

[DEBUG] Running docker inspect theucke-dbp-dbp
[DEBUG] Running docker run -d -it --name=theucke-dbp-dbp --hostname=stretch -v=/neteng/theucke/opx/dbp:/mnt -v=/home/theucke/.gitconfig:/etc/skel/.gitconfig:ro -e=UID=9438 -e=GID=3000 -e=TZ=US/Pacific-New -e=DEBFULLNAME=Dell EMC -e=DEBEMAIL=ops-dev@lists.openswitch.net -e=EXTRA_SOURCES=deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx
 opxhub/gbp:v1.0.5-stretch-dev bash -l
```
