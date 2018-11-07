from networkx import DiGraph

import dbp

MAKEFILE = dbp.MAKE_HEAD
MAKEFILE += """\t1 \\
\t2
1/${STAMP}: 2/${STAMP}
"""
MAKEFILE += dbp.MAKE_TAIL


def test_deb_build_options_string():
    assert dbp.deb_build_options_string(False) == "-e=DEB_BUILD_OPTIONS="
    assert (
        dbp.deb_build_options_string(True) == "-e=DEB_BUILD_OPTIONS=nostrip noopt debug"
    )
    assert dbp.deb_build_options_string(False, 5) == "-e=DEB_BUILD_OPTIONS=parallel=5"
    assert (
        dbp.deb_build_options_string(True, 5)
        == "-e=DEB_BUILD_OPTIONS=nostrip noopt debug parallel=5"
    )


def test_irun():
    assert dbp.irun(["true"], True) == 0
    assert dbp.irun(["true"], False) == 0
    assert dbp.irun(["false"], True) == 1
    assert dbp.irun(["false"], False) == 1
    assert dbp.irun([], True) == -1
    assert dbp.irun([], False) == -1


def test_makefile():
    g = DiGraph()
    g.add_node(1)
    g.add_node(2)
    g.add_edge(1, 2)
    assert dbp.makefile(g) == MAKEFILE
