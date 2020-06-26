package base

import "net"

func GetIPByHost(host string, defaultIP string) string {
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return defaultIP
	}

	return addr.String()
}
