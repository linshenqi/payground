package base

import (
	"net"
	"strings"
)

func GetIPByHost(host string, defaultIP string) string {
	targetHost := host
	vals := strings.Split(host, "://")
	if len(vals) > 1 {
		targetHost = strings.Split(vals[1], "/")[0]
	}

	addr, err := net.ResolveIPAddr("ip", targetHost)
	if err != nil {
		return defaultIP
	}

	return addr.String()
}
