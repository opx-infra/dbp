package main

import (
	"fmt"
	"log"
	"os"

	"github.com/opx-infra/dbp/workspace"
	"github.com/spf13/cobra"
)

const version = "18.12.1"

var (
	cnameFlag        string
	debugFlag        bool
	distFlag         string
	extraSourcesFlag string
	imageFlag        string
	pathFlag         string
	releaseFlag      string
	removeFirstFlag  bool
	verboseFlag      bool
	ws               *workspace.Workspace
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
	rootCmd := &cobra.Command{
		Use:     "dbp",
		Version: version,
		Short:   "Develop for Debian-based environments",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			ws, err = workspace.NewWorkspace(
				debugFlag,
				verboseFlag,
				pathFlag,
				cnameFlag,
				imageFlag,
				distFlag,
				releaseFlag,
				extraSourcesFlag,
			)
			if err != nil {
				log.Fatal(err)
			}

			if removeFirstFlag {
				cleanup(ws, true)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello world :)")
			fmt.Printf("cmd: %+v\n", cmd.PersistentFlags().Lookup("debug").Value.Type())
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cnameFlag, "cname", "", "", "container name")
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "", false, "build unstripped and unoptimized debug packages")
	rootCmd.PersistentFlags().StringVarP(&distFlag, "dist", "d", "stretch", "Debian distribution")
	rootCmd.PersistentFlags().StringVarP(&extraSourcesFlag, "extra-sources", "e", "DEFAULT", "extra apt sources")
	rootCmd.PersistentFlags().StringVarP(&imageFlag, "image", "i", "", "Docker image")
	rootCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", "", "workspace location")
	rootCmd.PersistentFlags().StringVarP(&releaseFlag, "release", "r", "unstable", "OPX release")
	rootCmd.PersistentFlags().BoolVarP(&removeFirstFlag, "rm-first", "", false, "manually remove container first")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "log debug messages to stderr")

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build Debian package(s) in a container",
		Run: func(cmd *cobra.Command, args []string) {
			alreadyRunning, err := ws.RunContainer(false)
			if err != nil {
				log.Fatal(err)
			}

			err = ws.BuildAllPackages(args)
			if err != nil {
				cleanup(ws, !alreadyRunning)
				log.Fatal(err)
			}

			cleanup(ws, !alreadyRunning)
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Ensure Docker container is present",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := ws.PullImage()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run persistent background container",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := ws.RunContainer(false)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	removeCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove persistent background container",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cleanup(ws, true)
		},
	}

	shellCmd := &cobra.Command{
		Use:   "shell",
		Short: "Launch interactive shell in container",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			alreadyRunning, err := ws.RunContainer(false)
			if err != nil {
				log.Fatal(err)
			}

			userCommand, _ := cmd.Flags().GetString("command")
			fullCommand := []string{"bash", "-l"}
			if userCommand != "" {
				fullCommand = append(fullCommand, "-c", userCommand)
			}
			ok := ws.DockerExec(fullCommand, "")

			cleanup(ws, !alreadyRunning)
			if !ok {
				log.Fatal("command exited with non-zero return code")
			}
		},
	}
	shellCmd.Flags().StringP("command", "c", "", "command to run")

	rootCmd.AddCommand(buildCmd)
	// TODO: makefile subcommand
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(shellCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
