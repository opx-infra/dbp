package main

import (
	"log"

	"github.com/opx-infra/dbp/workspace"
	flag "github.com/spf13/pflag"
)

func cleanup(ws *workspace.Workspace, remove bool) {
	if remove {
		err := ws.RemoveContainer()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// TODO: argument parsing and all that jazz
	// https://github.com/opx-infra/dbp/blob/15a41883cfdefc18f92b2e5a5885971c2a7d50d3/dbp/main.go
	// https://github.com/opx-infra/dbp/blob/master/dbp.py
	var cname = flag.StringP("cname", "", "", "container name")
	var debug = flag.BoolP("debug", "", false, "build unstripped and unoptimized debug packages")
	var dist = flag.StringP("dist", "d", "stretch", "Debian distribution")
	var extraSources = flag.StringP("extra-sources", "e", "DEFAULT", "extra apt sources")
	var image = flag.StringP("image", "i", "", "Docker image")
	var path = flag.StringP("path", "p", "", "workspace location")
	var release = flag.StringP("release", "r", "unstable", "OPX release")
	var removeFirst = flag.BoolP("rm-first", "", false, "manually remove container first")
	var verbose = flag.BoolP("verbose", "v", false, "log debug messages to stderr")
	flag.Parse()

	ws, err := workspace.NewWorkspace(
		*debug,
		*verbose,
		*path,
		*cname,
		*image,
		*dist,
		*release,
		*extraSources,
	)
	if err != nil {
		log.Fatal(err)
	}

	if *removeFirst {
		cleanup(ws, true)
	}

	alreadyRunning, err := ws.RunContainer(false)
	if err != nil {
		log.Fatal(err)
	}

	ok := ws.DockerExec([]string{"bash", "-l", "-c", "ls -l"}, "")
	if !ok {
		cleanup(ws, !alreadyRunning)
		log.Fatal("command exited with non-zero return code")
	}

	// err = ws.BuildPackage("dbp")
	// if err != nil {
	// cleanup(ws, !alreadyRunning)
	// log.Fatal(err)
	// }

	cleanup(ws, !alreadyRunning)
}
