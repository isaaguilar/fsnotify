//go:build !windows && !darwin
// +build !windows,!darwin

package internal

import (
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	SyscallEACCES = syscall.EACCES
	UnixEACCES    = unix.EACCES
)

var maxfiles uint64

// Go 1.19 will do this automatically: https://go-review.googlesource.com/c/go/+/393354/
func SetRlimit() {
	var l syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &l)
	if err == nil && l.Cur != l.Max {
		l.Cur = l.Max
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &l)
	}
	maxfiles = uint64(l.Cur)
}

func Maxfiles() uint64 { return maxfiles }
