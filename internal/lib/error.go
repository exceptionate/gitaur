package lib

import (
	"errors"
	"net"
	"net/http"
	"os"
	"syscall"
)

func GetError(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, http.ErrServerClosed) {
		return "server closed"
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			switch sysErr.Err {
			case syscall.EADDRINUSE:
				return "address already in use"
			case syscall.EACCES:
				return "permission denied"
			case syscall.EMFILE, syscall.ENFILE:
				return "too many open files"
			case syscall.EINVAL:
				return "invalid argument or port"
			case syscall.ENOBUFS, syscall.ENOMEM:
				return "system out of resources"
			}
		}
	}

	return err.Error()
}
