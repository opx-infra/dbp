package main

import (
	"log"

	"github.com/opx-infra/dbp/workspace"
	"github.com/spf13/cobra"
)

const version = "18.12.1"

type rootFlags struct {
	cname        string
	debug        bool
	dist         string
	extraSources string
	image        string
	path         string
	release      string
	removeFirst  bool
	verbose      bool
}

var (
	flags rootFlags
	ws    *workspace.Workspace
)

var rootCmd = &cobra.Command{
	Use:     "dbp",
	Version: version,
	Short:   "Develop for Debian-based environments",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		ws, err = workspace.NewWorkspace(
			flags.debug,
			flags.verbose,
			flags.path,
			flags.cname,
			flags.image,
			flags.dist,
			flags.release,
			flags.extraSources,
		)
		if err != nil {
			log.Fatal(err)
		}

		if flags.removeFirst {
			cleanup(ws, true)
		}
	},
}

var buildCmd = &cobra.Command{
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

var pullCmd = &cobra.Command{
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

var runCmd = &cobra.Command{
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

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove persistent background container",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cleanup(ws, true)
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Launch interactive shell in container",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		alreadyRunning, err := ws.RunContainer(false)
		if err != nil {
			log.Fatal(err)
		}

		userCommand, err := cmd.Flags().GetString("command")
		if err != nil {
			log.Fatal(err)
		}
		fullCommand := []string{"bash", "-l"}
		if userCommand != "" {
			fullCommand = append(fullCommand, "-c", userCommand)
		}

		ok := ws.DockerExec(fullCommand, "")

		cleanup(ws, !alreadyRunning)
		if !ok {
			log.Fatalf("command [%s] exited with non-zero return code\n", userCommand)
		}
	},
}

func init() {
	log.SetFlags(0)

	rootCmd.PersistentFlags().StringVarP(&flags.cname, "cname", "", "", "container name")
	rootCmd.PersistentFlags().BoolVarP(&flags.debug, "debug", "", false, "build unstripped and unoptimized debug packages")
	rootCmd.PersistentFlags().StringVarP(&flags.dist, "dist", "d", "stretch", "Debian distribution")
	rootCmd.PersistentFlags().StringVarP(&flags.extraSources, "extra-sources", "e", "DEFAULT", "extra apt sources")
	rootCmd.PersistentFlags().StringVarP(&flags.image, "image", "i", "", "Docker image")
	rootCmd.PersistentFlags().StringVarP(&flags.path, "path", "p", "", "workspace location")
	rootCmd.PersistentFlags().StringVarP(&flags.release, "release", "r", "unstable", "OPX release")
	rootCmd.PersistentFlags().BoolVarP(&flags.removeFirst, "rm-first", "", false, "manually remove container first")
	rootCmd.PersistentFlags().BoolVarP(&flags.verbose, "verbose", "v", false, "log debug messages to stderr")

	shellCmd.Flags().StringP("command", "c", "", "command to run")

	rootCmd.AddCommand(buildCmd)
	// TODO: makefile subcommand
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(shellCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func cleanup(ws *workspace.Workspace, remove bool) {
	if remove {
		err := ws.RemoveContainer()
		if err != nil {
			log.Fatal(err)
		}
	}
}
