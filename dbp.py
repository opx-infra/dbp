"""Simple program to manage gbp-docker container lifecycle."""

__version__ = "18.12.1"

import os
import shlex
import shutil
import subprocess
import sys
from pathlib import Path
from time import sleep
from typing import List

import click
import controlgraph
import docker
import networkx as nx

CONTEXT_SETTINGS = dict(help_option_names=["-h", "--help"])
IMAGE = "opxhub/gbp"
IMAGE_VERSION = "v2.0.5"
OPX_DEFAULT_SOURCES = "deb http://deb.openswitch.net/{} {} opx opx-non-free"


def info(s, nl=True):
    click.secho(s, nl=nl, fg="blue", err=True)


def error(s, nl=True):
    click.secho(s, nl=nl, fg="red", err=True)


class Workspace:
    """Manages the location of the workspace and its container."""

    def __init__(self, cname, debug, dist, extra_sources, image, release):
        self.debug = debug
        self.dist = dist
        self.image = image
        self.interactive = sys.stdin.isatty()
        self.release = release

        self.path = Path.cwd()
        if Path(self.path / ".git").is_dir():
            self.path = self.path.parent

        if cname:
            self.cname = cname
        else:
            self.cname = "{}-dbp-{}".format(os.getenv("USER"), self.path.name)

        # Sources order of preference
        # 1. extra_sources variable (set to "" for none)
        # 2. ./extra_sources.list file
        # 3. $HOME/.extra_sources.list file
        # 4. Default OPX sources
        if extra_sources == "DEFAULT":
            sources = [
                self.path / "extra_sources.list",
                Path.home() / ".extra_sources.list",
            ]
            for s in sources:
                if s.exists():
                    extra_sources = s.read_text()
                    self.extra_sources = extra_sources
                    break

        if extra_sources == "DEFAULT":
            self.extra_sources = OPX_DEFAULT_SOURCES.format(dist, release)
        else:
            self.extra_sources = extra_sources

        self.env = {
            "DEB_BUILD_OPTIONS": "nostrip noopt debug" if self.debug else "",
            "DEBEMAIL": os.getenv("DEBEMAIL", "ops-dev@lists.openswitch.net"),
            "DEBFULLNAME": os.getenv("DEBFULLNAME", "Dell EMC"),
            "EXTRA_SOURCES": self.extra_sources,
            "GID": os.getgid(),
            "TZ": "/".join(Path("/etc/localtime").resolve().parts[-2:]),
            "UID": os.getuid(),
        }

        self.volumes = {str(self.path): {"bind": "/mnt", "mode": "rw"}}
        gitconfig = Path(Path.home() / ".gitconfig")
        if gitconfig.exists():
            self.volumes[str(gitconfig)] = {
                "bind": "/etc/skel/.gitconfig",
                "mode": "ro",
            }

        self.client = docker.from_env()

    def builddepends_graph(self) -> nx.DiGraph:
        dirs = [p for p in self.path.iterdir() if p.is_dir()]
        return controlgraph.graph(controlgraph.parse_all_controlfiles(dirs))

    def buildpackage(self, directory: Path, gbp_options: str) -> int:
        """Runs gbp buildpackage (or debuild for 3.0 (git) packages) in the container."""
        pkg_format = "1.0"
        if Path(directory / "debian/source/format").exists():
            pkg_format = Path(directory / "debian/source/format").read_text().strip()

        if "3.0 (git)" in pkg_format:
            build_cmd = ["debuild"]
        else:
            build_cmd = ["gbp", "buildpackage"]

        if gbp_options and build_cmd[0] == "gbp":
            build_cmd.extend(shlex.split(gbp_options))

        return self.docker_exec(build_cmd, "/mnt/{}".format(directory))

    def docker_exec(self, command: List[str], workdir="/mnt") -> int:
        """Low level function to execute a command in an already running container."""
        if not self.container:
            return 1

        full_cmd = [
            "docker",
            "exec",
            "-it" if self.interactive else "-t",
            "--workdir={}".format(workdir),
        ]
        for k, v in self.env.items():
            full_cmd.append("-e={}={}".format(k, v))
        full_cmd.append(self.cname)
        full_cmd.append("/entrypoint.sh")
        full_cmd.extend(command)

        proc = subprocess.run(
            full_cmd, stdin=sys.stdin, stdout=sys.stdout, stderr=sys.stderr
        )
        return proc.returncode

    def docker_remove(self) -> int:
        """Low level function to remove our running container."""
        try:
            container = self.client.containers.get(self.cname)
        except docker.errors.NotFound:
            return

        info("Removing container {}...".format(self.cname), nl=False)
        container.remove(force=True)
        info("Done!")

    def docker_run(self) -> bool:
        """Runs the container and returns True if it didn't exist before."""
        containers = self.client.containers.list(filters={"name": self.cname})
        if containers:
            self.container = containers[0]
            return False

        info("Starting container {}...".format(self.cname), nl=False)
        self.container = self.client.containers.run(
            image=self.image,
            detach=True,
            auto_remove=False,
            environment=self.env,
            hostname=self.dist,
            init=True,
            name=self.cname,
            remove=False,
            stdin_open=sys.stdin.isatty(),
            tty=True,
            volumes=self.volumes,
            entrypoint="bash",
            command=[],
        )
        info("Done!")

        return True


def generate_makefile(g: nx.DiGraph) -> str:
    bob = ""  # the builder
    dep_lines = []  # makefile dependency lines to print

    bob += """.PHONY: all
all:
ALL_REPOS = \\
"""
    nodes = sorted([n for n in nx.dfs_postorder_nodes(g)])
    for n in nodes:
        if n == nodes[len(nodes) - 1]:
            bob += "\t{n}\n".format(n=n)
        else:
            bob += "\t{n} \\\n".format(n=n)
        for dep in g.successors(n):
            dep_lines.append(".{}.stamp: .{}.stamp\n".format(n, dep))
    for line in dep_lines:
        bob += line

    bob += """ALL_REPO_STAMPS := $(patsubst %,.%.stamp,${ALL_REPOS})
TIMESTAMP = $(shell date '+%F %T')

all: ${ALL_REPO_STAMPS}

${ALL_REPO_STAMPS}: REPO = $(@:.%.stamp=%)
${ALL_REPO_STAMPS}: LOG = ${REPO}.log
${ALL_REPO_STAMPS}:
\t@echo ${TIMESTAMP} Starting dbp build ${REPO}
\t@dbp --cname "$${USER}-dbp-parallel-${REPO}" build ${REPO} >${LOG} 2>&1
\t@: >$@"""

    return bob


# Command line interface ##############################################################

pass_workspace = click.make_pass_decorator(Workspace)


@click.group(context_settings=CONTEXT_SETTINGS)
@click.option("--cname", envvar="CNAME", metavar="NAME", help="Custom container name.")
@click.option("--debug", is_flag=True, help="Set nostrip, noopt, debug.")
@click.option(
    "--dist",
    "-d",
    envvar="DIST",
    default="stretch",
    metavar="DIST",
    help="Debian distribution.",
)
@click.option(
    "--extra-sources",
    "-e",
    envvar="EXTRA_SOURCES",
    default="DEFAULT",
    metavar="SOURCES",
    help="Extra Apt sources to add along with OPX.",
)
@click.option(
    "--image",
    "-i",
    envvar="IMAGE",
    default=IMAGE,
    metavar="IMAGE",
    help="Docker image.",
)
@click.option(
    "--release",
    "-r",
    envvar="RELEASE",
    default="unstable",
    metavar="RELEASE",
    help="OPX release.",
)
@click.option(
    "--rm-first/--no-rm-first",
    default=False,
    help="Remove any active container before running.",
)
@click.version_option(__version__)
@click.pass_context
def cli(ctx, cname, debug, dist, extra_sources, image, release, rm_first):
    """Build Debian packages in a managed Debian development environment."""
    ctx.obj = Workspace(cname, debug, dist, extra_sources, image, release)

    ws = ctx.obj
    if ":" not in ws.image:
        ws.image = "{}:{}-{}".format(ws.image, IMAGE_VERSION, ws.dist)
        if ctx.invoked_subcommand != "build":
            ws.image += "-dev"

    # check for prereqs
    if shutil.which("docker") is None:
        error("Docker not found in PATH. Please install docker.")
        sys.exit(1)
    if shutil.which("git") is None:
        error("Git not found in PATH. Please install git.")
        sys.exit(1)

    # ensure Docker image is present
    if ctx.invoked_subcommand in ["run", "build", "pull", "shell"]:
        try:
            ws.client.images.get(ws.image)
        except docker.errors.ImageNotFound:
            info("Pulling image {}...".format(ws.image))
            ws.client.images.pull(ws.image)

    if rm_first:
        ws.docker_remove()


@cli.command()
@click.option(
    "--gbp",
    "-g",
    envvar="GBP_OPTS",
    default="",
    metavar="OPTS",
    help="Options to pass to Git-buildpackage.",
)
@click.option("--print-targets/--no-print-targets", help="Print build order and exit.")
@click.argument("targets", nargs=-1)
@pass_workspace
def build(ws, gbp, print_targets, targets):
    """Build packages.

    Builds targets specified. Otherwise builds a single repo if run from within the repo
    or builds all repos with d/control files in dependency order.
    """
    if not targets and ws.path == Path.cwd():
        # If run from workspace root with no targets, build all in dependency order
        targets = tuple(nx.dfs_postorder_nodes(ws.builddepends_graph()))
    elif not targets:
        # If run from a directory in the workspace with no targets, build the directory
        targets = (Path.cwd().stem,)

    info("Building {} repositories: {}".format(len(targets), " ".join(targets)))
    if print_targets:
        sys.exit(0)

    rc = 0
    remove_container = ws.docker_run()

    for t in targets:
        info("--- Building {}...".format(t))
        rc = ws.buildpackage(Path(t), gbp)
        if rc:
            error("Building {} failed with return code {}.".format(t, rc))
            break

    if remove_container:
        ws.docker_remove()

    sys.exit(rc)


@cli.command()
@pass_workspace
def makefile(ws):
    """Write a Makefile suitable for parallel building."""
    g = ws.builddepends_graph()
    sys.stdout.write(generate_makefile(g))
    sys.stdout.flush()


@cli.command()
@pass_workspace
def pull(ws):
    """Ensure Docker image is present.

    Pulling is done when creating a Workspace object, so this isn't necessary.
    """


@cli.command()
@pass_workspace
def run(ws):
    """Start a persistent workspace container."""
    ws.docker_run()


@cli.command()
@pass_workspace
def rm(ws):
    """Remove any workspace containers."""
    ws.docker_remove()


@cli.command()
@click.option("--command", "-c", metavar="CMD", help="Non-interactive command to run.")
@pass_workspace
def shell(ws, command):
    """Execute a shell or a command in a workspace container."""
    full_cmd = ["bash", "-l"]
    if command:
        full_cmd.extend(["-c", command])

    remove_container = ws.docker_run()

    rc = ws.docker_exec(full_cmd)

    if remove_container:
        ws.docker_remove()

    sys.exit(rc)
