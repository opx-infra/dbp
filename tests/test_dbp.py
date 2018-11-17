from networkx import DiGraph

import dbp

MAKEFILE = """.PHONY: all
all:
ALL_REPOS = \\
	1 \\
	2
.1.stamp: .2.stamp
ALL_REPO_STAMPS := $(patsubst %,.%.stamp,${ALL_REPOS})
TIMESTAMP = $(shell date '+%F %T')

all: ${ALL_REPO_STAMPS}

${ALL_REPO_STAMPS}: REPO = $(@:.%.stamp=%)
${ALL_REPO_STAMPS}: LOG = ${REPO}.log
${ALL_REPO_STAMPS}:
\t@echo ${TIMESTAMP} Starting dbp build ${REPO}
\t@dbp --cname "$${USER}-dbp-parallel-${REPO}" build ${REPO} >${LOG} 2>&1
\t@: >$@"""


def test_makefile():
    g = DiGraph()
    g.add_node(1)
    g.add_node(2)
    g.add_edge(1, 2)
    assert dbp.generate_makefile(g) == MAKEFILE
