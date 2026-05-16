//go:build windows

package platform

import (
	"net"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func DialIPC(path string) (net.Conn, error) {
	timeout := 10 * time.Second
	return winio.DialPipe(path, &timeout)
}
