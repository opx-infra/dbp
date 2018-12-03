package main

import (
	"log"

	"github.com/opx-infra/dbp/workspace"
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
	ws, err := workspace.NewWorkspace(
		false,
		true,
		"",
		"",
		"opxhub/gbp:v2-stretch-dev",
		"stretch",
		"unstable",
		"DEFAULT",
	)
	if err != nil {
		log.Fatal(err)
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

	// err = ws.BuildPackage("dbp-go")
	// if err != nil {
	// log.Fatal(err)
	// }

	cleanup(ws, !alreadyRunning)
}
