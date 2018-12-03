package workspace

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	defaultDist    = "stretch"
	defaultRelease = "unstable"
	dockerImage    = "opxhub/gbp:v2.0.5"
)

// Workspace represents a Docker container and directory
type Workspace struct {
	CName        string            // Container name
	Client       *client.Client    // Docker client
	Debug        bool              // Build debug, unoptimized, unstripped packages
	Dist         string            // Debian distribution to build against
	Env          map[string]string // Environment variables for the container
	ExtraSources string            // Extra apt sources to use
	Image        string            // Docker image in use
	Interactive  bool              // True if Stdin is a terminal (see golang.org/x/crypto/ssh/terminal)
	DebugLogger  *log.Logger       // Debug logger
	InfoLogger   *log.Logger       // Info logger
	Path         string            // Location of workspace
	Release      string            // OPX release to build against
	Volumes      []string          // Volume mount mapping for the container
}

// NewWorkspace constructs a workspace to run package builds in
func NewWorkspace(debug, verbose bool, path, cname, image, dist, release, extraSources string) (*Workspace, error) {
	var ws Workspace
	var err error

	var out io.Writer
	if verbose {
		out = os.Stderr
	} else {
		out = ioutil.Discard
	}
	ws.DebugLogger = log.New(out, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)
	ws.InfoLogger = log.New(os.Stderr, "", 0)

	// Process arguments
	// debug
	ws.Debug = debug
	var debugOptions string
	if debug {
		debugOptions = "nostrip noopt debug"
	}

	// path
	ws.Path, err = filepath.Abs(path)
	if err != nil {
		return &ws, errors.Wrap(err, "pwd failed")
	}
	gitDir := filepath.Join(ws.Path, ".git")
	_, err = os.Stat(gitDir)
	if err == nil {
		// .git/ exists, use parent
		ws.Path = filepath.Dir(ws.Path)
	} else if os.IsNotExist(err) {
		// .git/ does not exist, use cwd as-is
	} else {
		return &ws, errors.Wrap(err, "stat on .git failed")
	}

	// cname
	ws.CName = cname
	if ws.CName == "" {
		ws.CName = fmt.Sprintf("%s-dbp-%s", os.Getenv("USER"), filepath.Base(ws.Path))
	}

	// dist
	if dist == "" {
		dist = defaultDist
	}
	ws.Dist = dist

	// release
	if release == "" {
		release = defaultRelease
	}
	ws.Release = release

	// image
	if image == "" {
		image = fmt.Sprintf("%s-%s-dev", dockerImage, ws.Dist)
	}
	ws.Image = image

	// extraSources
	// Sources order of preference
	// 1. extraSources variable (set to "" for none)
	// 2. ./extra_sources.list file
	// 3. $HOME/.extra_sources.list file
	// 4. Default OPX sources
	err = ws.setExtraSources([]string{
		"./extra_sources.list",
		filepath.Join(getenv("HOME", "/"), ".extra_sources.list"),
	}, extraSources)
	if err != nil {
		return &ws, errors.Wrap(err, "setting extra sources")
	}

	// Interactive
	_, err = unix.IoctlGetTermios(int(os.Stdin.Fd()), ioctlReadTermios)
	ws.Interactive = err == nil

	// Set up environment
	// For the timezone, setting $TZ is more portable than mounting /etc/localtime
	tz, err := filepath.EvalSymlinks("/etc/localtime")
	if err != nil {
		return &ws, errors.Wrap(err, "resolving localtime symlink")
	}
	// We only want the last two parts, e.g. America/LosAngeles
	if tzParts := strings.Split(tz, "/"); len(tzParts) > 1 {
		tz = fmt.Sprintf("%s/%s", tzParts[len(tzParts)-2], tzParts[len(tzParts)-1])
	}
	ws.Env = map[string]string{
		"DEB_BUILD_OPTIONS": debugOptions,
		"DEBEMAIL":          getenv("DEBEMAIL", "ops-dev@lists.openswitch.net"),
		"DEBFULLNAME":       getenv("DEBFULLNAME", "Dell EMC"),
		"EXTRA_SOURCES":     ws.ExtraSources,
		"GID":               strconv.Itoa(os.Getgid()),
		"TZ":                tz,
		"UID":               strconv.Itoa(os.Getuid()),
	}

	// Set up volume mounts
	ws.Volumes = []string{fmt.Sprintf("%s:/mnt:rw", ws.Path)}
	gitconfig := filepath.Join(getenv("HOME", "."), ".gitconfig")
	_, err = os.Stat(gitconfig)
	if err == nil {
		// ~/.gitconfig exists, mount in container
		ws.Volumes = append(ws.Volumes, fmt.Sprintf("%s:/etc/skel/.gitconfig:ro", gitconfig))
	} else if os.IsNotExist(err) {
		// ~/.gitconfig does not exist
	} else {
		return &ws, errors.Wrap(err, "stat on gitconfig failed")
	}

	ws.DebugLogger.Printf("Created workspace: %+v\n", ws)
	return &ws, nil
}
