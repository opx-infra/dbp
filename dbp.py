"""Simple program to manage gbp-docker container lifecycle."""

__version__ = "0.4.0~pre"

import argparse
import logging
import os
import shutil
import sys

from collections import namedtuple
from pathlib import Path
from queue import Queue
from subprocess import run, PIPE, DEVNULL, STDOUT
from threading import Thread
from typing import List

import networkx as nx

PWD = Path.cwd()
UID = os.getuid()
GID = os.getgid()
USER = os.getenv("USER")

CONTAINER_NAME = "{}-dbp-{}".format(USER, PWD.stem)
IMAGE = "opxhub/gbp"
IMAGE_VERSION = "v1.0.0"


L = logging.getLogger("dbp")
L.addHandler(logging.NullHandler())


def image_name(image: str, dist: str, dev: bool) -> str:
    """Returns the Docker image to use, allowing for custom images."""
    if ":" in image:
        return image

    if dev:
        template = "{}:{}-{}-dev"
    else:
        template = "{}:{}-{}"

    return template.format(image, IMAGE_VERSION, dist)


def irun(cmd: List[str], quiet=False, file=None) -> int:
    """irun runs an interactive command."""
    L.debug("Running {}".format(" ".join(cmd)))
    if file is not None:
        proc = run(cmd, stdout=file.open("w"), stderr=STDOUT)
    elif quiet:
        proc = run(cmd, stdout=DEVNULL, stderr=DEVNULL)
    else:
        proc = run(cmd, stdin=sys.stdin, stdout=sys.stdout, stderr=sys.stderr)
    return proc.returncode


def container_exists() -> bool:
    """Returns true if our dbp container can be inspected"""
    return irun(["docker", "inspect", CONTAINER_NAME], quiet=True) == 0


def container_running(dist: str) -> bool:
    """Returns true if our dbp container is running"""
    proc = run(
        ["docker", "inspect", CONTAINER_NAME, "--format={{.State.Running}}"],
        stdout=PIPE,
        stderr=DEVNULL,
    )
    return proc.returncode == 0 and "true" in str(proc.stdout)


def buildpackage(
    dist: str, target: Path, sources: str, gbp_options: str, quiet=False, file=None
) -> int:
    """Runs gbp buildpackage --git-export-dir=pool/{dist}-amd64/{target}

    Container must already be started.
    """
    if not target.exists():
        L.error("Build target `{}` does not exist".format(target))
        return 1

    cmd = ["docker", "exec", "-t"]
    if sys.stdin.isatty():
        cmd.append("-i")
    cmd.extend(
        [
            "--user=build",
            "-e=UID={}".format(UID),
            "-e=GID={}".format(GID),
            "-e=EXTRA_SOURCES={}".format(sources),
            CONTAINER_NAME,
            "build",
            target.stem,
            gbp_options,
        ]
    )

    return irun(cmd, quiet, file)


def docker_run(image: str, dist: str, sources: str, dev=True) -> int:
    if container_exists():
        L.warning("Container already exists.")
        return 0

    cmd = [
        "docker",
        "run",
        "-d",
        "-it",
        "--name={}".format(CONTAINER_NAME),
        "--hostname={}".format(dist),
        "-v={}:/mnt".format(PWD),
        "-v={}/.gitconfig:/etc/skel/.gitconfig:ro".format(Path.home()),
        "-e=UID={}".format(UID),
        "-e=GID={}".format(GID),
        "-e=EXTRA_SOURCES={}".format(sources),
        image_name(image, dist, dev),
    ]

    if not dev:
        cmd = cmd + ["bash", "-l"]

    return irun(cmd, quiet=True)


def docker_shell(image: str, dist: str, sources: str, command=None) -> int:
    """Runs bash in a development container

    Uses (and starts if necessary) any matching pre-existing dbp container.
    """
    if container_exists():
        cmd = [
            "docker",
            "exec",
            "-it",
            "--user=build",
            "-e=UID={}".format(UID),
            "-e=GID={}".format(GID),
            "-e=EXTRA_SOURCES={}".format(sources),
            CONTAINER_NAME,
            "bash",
            "-l",
        ]
        if not container_running(dist):
            docker_start(dist)
    else:
        cmd = [
            "docker",
            "run",
            "--rm",
            "-it",
            "--name={}".format(CONTAINER_NAME),
            "--hostname={}".format(dist),
            "-v={}:/mnt".format(PWD),
            "-v={}/.gitconfig:/etc/skel/.gitconfig:ro".format(Path.home()),
            "-e=UID={}".format(UID),
            "-e=GID={}".format(GID),
            "-e=EXTRA_SOURCES={}".format(sources),
            image_name(image, dist, True),
        ]

    if command is not None:
        cmd = cmd + ["bash", "-l", "-c", command]

    return irun(cmd, quiet=False)


def docker_start(dist: str) -> int:
    """Runs docker start and returns the return code"""
    cmd = ["docker", "start", CONTAINER_NAME]
    return irun(cmd, quiet=True)


def pull_images(image: str, dist: str) -> int:
    """Runs docker pull for both build and development images and returns the return code"""
    cmd = ["docker", "pull", image_name(image, dist, False)]
    rc = irun(cmd)
    if rc != 0:
        return rc

    cmd = ["docker", "pull", image_name(image, dist, True)]
    return irun(cmd)


def remove_container() -> int:
    """Runs docker rm -f for the dbp container"""
    if container_exists():
        cmd = ["docker", "rm", "-f", CONTAINER_NAME]
        return irun(cmd, quiet=True)

    L.warning("Container does not exist.")
    return 1


WorkItem = namedtuple(
    "WorkItem",
    ["graph", "node", "built", "level", "dist", "sources", "options", "quiet"],
)


def build_worker(q):
    while True:
        work = q.get()
        L.debug("{}{} started".format("-" * work.level, work.node))

        if not work.built[work.node] and all(
            [work.built[n] for n in work.graph.predecessors(work.node)]
        ):
            L.info("{}{} building".format("-" * work.level, work.node))

            file = None
            if work.quiet:
                file = Path("{}-{}.build".format(work.level, work.node))
            rc = buildpackage(
                work.dist, Path(work.node), work.sources, work.options, work.quiet, file
            )

            if rc != 0:
                L.error("{}{} failure".format("-" * work.level, work.node))
            else:
                L.info("{}{} success!".format("-" * work.level, work.node))
                work.built[work.node] = True
                for n in list(work.graph.successors(work.node)):
                    if all([work.built[m] for m in work.graph.predecessors(n)]):
                        q.put(
                            WorkItem(
                                graph=work.graph,
                                node=n,
                                built=work.built,
                                level=work.level + 1,
                                dist=work.dist,
                                sources=work.sources,
                                options=work.options,
                                quiet=work.quiet,
                            )
                        )

        L.debug("{}{} finished".format("-" * work.level, work.node))
        q.task_done()


def cmd_build(args: argparse.Namespace) -> int:
    """Run with dbp build"""
    rc = 0
    remove = True

    if container_exists():
        remove = False
    else:
        rc = docker_run(args.image, args.dist, args.extra_sources, dev=False)
        if rc != 0:
            L.error("Could not run container")
            return rc

    if not container_running(args.dist):
        rc = docker_start(args.dist)
        if rc != 0:
            L.error("Could not start stopped container")
            return 1

    G = nx.drawing.nx_pydot.read_dot(sys.stdin)
    built = {n: False for n in G}

    no_deps = [n for n in G if G.in_degree(n) == 0]
    q = Queue(maxsize=len(no_deps))
    threads = []
    for i in range(args.parallel):
        worker = Thread(target=build_worker, args=(q,))
        worker.setDaemon(True)
        worker.start()
        threads.append(worker)

    quiet = args.parallel > 1
    for n in no_deps:
        q.put(WorkItem(G, n, built, 0, args.dist, args.extra_sources, args.gbp, quiet))

    try:
        q.join()
    except KeyboardInterrupt:
        L.info("Got SIGINT")
        q.queue.clear()
    finally:
        if remove:
            remove_container()

    return 0


def cmd_pull(args: argparse.Namespace) -> int:
    """Run with dbp pull"""
    return pull_images(args.image, args.dist)


def cmd_rm(args: argparse.Namespace) -> int:
    """Run with dbp rm"""
    remove_container()
    return 0


def cmd_run(args: argparse.Namespace) -> int:
    """Run with dbp run"""
    return docker_run(args.image, args.dist, args.extra_sources, dev=True)


def cmd_shell(args: argparse.Namespace) -> int:
    """Run with dbp shell"""
    return docker_shell(args.image, args.dist, args.extra_sources, args.command)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)

    # general arguments
    parser.add_argument(
        "--version", "-V", action="store_true", help="print program version"
    )
    parser.add_argument(
        "--verbose", "-v", help="-v for info, -vv for debug", action="count", default=0
    )
    parser.add_argument(
        "--dist", "-d", help="Debian distribution", default=os.getenv("DIST", "stretch")
    )
    parser.add_argument(
        "--extra-sources",
        "-e",
        help="Apt-style sources",
        default=os.getenv("EXTRA_SOURCES", ""),
    )
    parser.add_argument("--image", "-i", help="Docker image to use", default=IMAGE)

    sps = parser.add_subparsers(help="commands")

    # build subcommand
    build_parser = sps.add_parser("build", help="run gbp buildpackage")
    build_parser.add_argument(
        "targets", nargs="*", type=Path, help="directories to build"
    )
    build_parser.add_argument(
        "--gbp", "-g", default="", help="additional git-buildpackage options to pass"
    )
    build_parser.add_argument(
        "--parallel",
        "-p",
        default=1,
        type=int,
        help="build in parallel; must provide reverse controlgraph on stdin",
    )
    build_parser.set_defaults(func=cmd_build)

    # pull subcommand
    pull_parser = sps.add_parser("pull", help="pull latest images")
    pull_parser.set_defaults(func=cmd_pull)

    # rm subcommand
    rm_parser = sps.add_parser("rm", help="remove workspace container")
    rm_parser.set_defaults(func=cmd_rm)

    # run subcommand
    run_parser = sps.add_parser("run", help="run development container in background")
    run_parser.set_defaults(func=cmd_run)

    # shell subcommand
    shell_parser = sps.add_parser("shell", help="launch development environment")
    shell_parser.add_argument("--command", "-c", help="command to run noninteractively")
    shell_parser.set_defaults(func=cmd_shell)

    args = parser.parse_args()

    if args.version:
        print("dbp {}".format(__version__))
        return 0

    # set up logging
    logging.basicConfig(
        format="[%(levelname)s] %(message)s", level=10 * (3 - min(args.verbose, 2))
    )

    # check for prereqs
    if shutil.which("docker") is None:
        L.error("Docker not found in PATH. Please install docker.")
        sys.exit(1)

    # read sources from ./extra_sources and ~/.extra_sources
    if args.extra_sources == "":
        extra_sources = [Path("extra_sources"), Path.home() / ".extra_sources"]
        for s in extra_sources:
            if s.exists():
                args.extra_sources = s.read_text()
                break

    if args.extra_sources != "":
        L.info("Loaded extra sources:\n{}".format(args.extra_sources))

    # print help if no subcommand specified
    if getattr(args, "func", None) is None:
        parser.print_help()
        return 0

    # run subcommand
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
