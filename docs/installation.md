# Installation

## Prerequisites

1. Python 3.5+
1. Docker

## Install

```bash
$ pip3 install --upgrade dbp
```

To avoid installing into the global namespace, install into a virtualenv.

```bash
$ python3 -mvenv ~/.dbp && ~/.dbp/bin/pip install dbp
$ export PATH=$PATH:$HOME/.dbp/bin
```
