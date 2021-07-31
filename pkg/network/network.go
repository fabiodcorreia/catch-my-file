package network

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Socket Type.
const Type = `tcp4`

// Hostname return the marchine hostname.
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("peer-%d", time.Now().Unix())
	}
	return hostname
}

// GetLocalIP returns the local ip besides the localhost.
func GetLocalIP() (*net.IPNet, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet, nil
			}
		}
	}
	return nil, fmt.Errorf("network get local ip found")
}

// SameNetworks return if the ip provided is part tof the same network
// of the ipnet.
func SameNetwork(ipnet *net.IPNet, ip net.IP) bool {
	return ipnet.Contains(ip)
}
