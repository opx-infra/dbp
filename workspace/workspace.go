package workspace

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// TODO: builddepends graph (adjacency list?) in new package
// https://github.com/opx-infra/builddepends/blob/master/builddepends.go
// TODO: makefile from graph

// BuildAllPackages builds all packages specified in dependency order
func (ws *Workspace) BuildAllPackages(paths []string) error {
	if len(paths) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if cwd == ws.Path {
			// TODO: build dependency graph
			paths = append(paths, filepath.Base(cwd))
		} else {
			paths = append(paths, filepath.Base(cwd))
		}
	}

	for _, path := range paths {
		err := ws.BuildPackage(path)
		if err != nil {
			return errors.Wrap(err, path)
		}
	}
	return nil
}

// BuildPackage runs gbp buildpackage or debuild in a running container
func (ws *Workspace) BuildPackage(path string) error {
	if strings.HasPrefix(path, "/") {
		return errors.New("path must be relative to workspace")
	}
	if _, err := os.Stat(filepath.Join(ws.Path, path)); os.IsNotExist(err) {
		return errors.Errorf("path %s does not exist in workspace", path)
	}
	pkgFormat := "1.0"

	var pkgFormatFile string
	if strings.HasPrefix(path, "/") {
		pkgFormatFile = filepath.Join(path, "debian/source/format")
	} else {
		pkgFormatFile = filepath.Join(ws.Path, path, "debian/source/format")
	}
	data, err := ioutil.ReadFile(pkgFormatFile)
	if err != nil {
		if os.IsNotExist(err) {
			// That's fine, use default value
		} else {
			return errors.Wrap(err, "reading debian source format failed")
		}
	}
	pkgFormat = strings.Trim(string(data), "\n \t")

	var buildCmd []string
	switch pkgFormat {
	case "3.0 (git)":
		buildCmd = append(buildCmd, "debuild")
	case "1.0":
		fallthrough
	case "3.0 (native)":
		fallthrough
	case "3.0 (quilt)":
		fallthrough
	default:
		buildCmd = append(buildCmd, "gbp", "buildpackage")
	}

	ok := ws.DockerExec(buildCmd, fmt.Sprintf("/mnt/%s", path))
	if ok {
		return nil
	}

	return errors.New("package build failed")
}

// DockerExec runs a command in a running gbp-docker container
func (ws *Workspace) DockerExec(cmd []string, workDir string) bool {
	// TODO: use docker sdk
	// https://github.com/docker/cli/blob/master/cli/command/container/exec.go
	// https://godoc.org/github.com/docker/docker/api/types
	// https://godoc.org/github.com/docker/docker/client
	if workDir == "" {
		workDir = "/mnt"
	}
	fullCmd := []string{"exec", "--tty", fmt.Sprintf("--workdir=%s", workDir)}
	if ws.Interactive {
		fullCmd = append(fullCmd, "--interactive")
	}
	fullCmd = append(fullCmd, ws.envStrings(true)...)
	fullCmd = append(fullCmd, ws.CName, "/entrypoint.sh")
	fullCmd = append(fullCmd, cmd...)

	proc := exec.Command("docker", fullCmd...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	ws.DebugLogger.Printf("Running %v in %s...\n", cmd, workDir)
	proc.Run()
	return proc.ProcessState.Success()
}

// PullImage pulls our workspace Docker image
func (ws *Workspace) PullImage() error {
	_, err := ws.dockerClient()
	if err != nil {
		return err
	}
	ws.InfoLogger.Printf("Pulling image %s...\n", ws.Image)
	out, err := ws.Client.ImagePull(context.Background(), ws.Image, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "pulling image")
	}
	defer out.Close()
	// io.Copy(os.Stdout, out)
	if _, err := ioutil.ReadAll(out); err != nil {
		return errors.Wrap(err, "pulling image")
	}
	return nil
}

// RemoveContainer forcefully removes the workspace container
func (ws *Workspace) RemoveContainer() error {
	_, err := ws.dockerClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	containers, err := ws.Client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "listing containers")
	}
	for _, container := range containers {
		if len(container.Names) == 1 && container.Names[0] == fmt.Sprintf("/%s", ws.CName) {
			// Our container already exists!
			ws.DebugLogger.Printf("Removing container %s...\n", ws.CName)
			err = ws.Client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
				Force: true,
			})
			if err != nil {
				return errors.Wrap(err, "removing container")
			}
			return nil
		}
	}

	return nil
}

// RunContainer runs the Docker container in the background
func (ws *Workspace) RunContainer(pull bool) (bool, error) {
	_, err := ws.dockerClient()
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	// Check if container already exists and just return
	containers, err := ws.Client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, errors.Wrap(err, "listing containers")
	}
	for _, container := range containers {
		if len(container.Names) == 1 && container.Names[0] == fmt.Sprintf("/%s", ws.CName) {
			// Our container already exists!
			return true, nil
		}
	}

	// Check if image exists and pull if not
	haveImage := false
	images, err := ws.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return false, errors.Wrap(err, "listing images")
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == ws.Image {
				haveImage = true
				break
			}
		}
		if haveImage {
			break
		}
	}
	if pull || !haveImage {
		err = ws.PullImage()
		if err != nil {
			return false, err
		}
	}

	// Create container
	ws.DebugLogger.Printf("Running container %s...\n", ws.CName)
	res, err := ws.Client.ContainerCreate(ctx, &container.Config{
		Hostname:   ws.Dist,
		Tty:        true,
		OpenStdin:  ws.Interactive,
		Env:        ws.envStrings(false),
		Cmd:        []string{"-l"},
		Image:      ws.Image,
		Entrypoint: []string{"bash"},
	}, &container.HostConfig{
		Binds: ws.Volumes,
	}, nil, ws.CName)
	if err != nil {
		return false, errors.Wrap(err, "creating container")
	}

	// Start container
	ws.DebugLogger.Printf("Starting container %s...\n", ws.CName)
	if err := ws.Client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{}); err != nil {
		return false, errors.Wrap(err, "starting container")
	}

	return false, nil
}

// dockerClient returns the Workspace Docker client
func (ws *Workspace) dockerClient() (*client.Client, error) {
	if ws.Client == nil {
		var err error
		ws.Client, err = client.NewEnvClient()
		if err != nil {
			return nil, errors.Wrap(err, "getting docker client")
		}
	}
	return ws.Client, nil
}

// envStrings returns docker-ready env var pairs
func (ws *Workspace) envStrings(withOpt bool) []string {
	var opt string
	if withOpt {
		opt = "-e="
	}
	env := make([]string, len(ws.Env))
	i := 0
	for k, v := range ws.Env {
		env[i] = fmt.Sprintf("%s%s=%s", opt, k, v)
		i++
	}
	return env
}

// getenv is os.getenv with a default value for convenience
func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// setExtraSources reads from several locations and returns the appropriate apt source list
func (ws *Workspace) setExtraSources(files []string, extraSources string) error {
	if extraSources != "DEFAULT" {
		ws.ExtraSources = extraSources
	} else {
		// search list of files until one exists
		for _, f := range files {
			data, err := ioutil.ReadFile(f)
			if err != nil {
				// Can't read? We don't care, move on
				continue
			}
			ws.ExtraSources = string(data)
			break
		}
	}

	if extraSources == "DEFAULT" {
		// no files were found, set to default value
		ws.ExtraSources = fmt.Sprintf(
			"deb http://deb.openswitch.net/%s %s opx opx-non-free",
			ws.Dist,
			ws.Release,
		)
	}

	return nil
}
