// +build darwin dragonfly freebsd netbsd openbsd

package workspace

import "golang.org/x/sys/unix"

// https://www.freebsd.org/cgi/man.cgi?query=tty&sektion=4<Paste>
// tcgetattr()
const ioctlReadTermios = unix.TIOCGETA
