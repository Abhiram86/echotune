//go:build !windows

package platform

import (
	"net"
	"time"
)

func DialIPC(path string) (net.Conn, error) {
	return net.DialTimeout("unix", path, 10*time.Second)
}
