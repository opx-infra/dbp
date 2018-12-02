package workspace

import "golang.org/x/sys/unix"

// http://man7.org/linux/man-pages/man4/tty_ioctl.4.html
// tcgetattr()
const ioctlReadTermios = unix.TCGETS
