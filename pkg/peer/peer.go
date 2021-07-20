// Peer package is containes all the componets related with the network peers.
//
// Components:
//
// - peerserver.go -> Registration and peer discovery
//
package peer

import (
	"net"
	"strings"
)

const (
	serviceName   = `_catchmyfile._tcp`
	serviceDomain = `local.`
)

// Peer defines a network peer that can send and receive files.
type Peer struct {
	// Name is the peers name.
	Name string
	// Address - os the net.IP address of the peer.
	Address net.IP
	// Port is the network port where the peer will receive connections.
	Port int
}

// newPeer will create a new instance of the struct Peer and return it.
//
// It receives a string with the name of the peer, the ip address and
// the port where the peer will receive connections for file transfers.
func newPeer(name string, address net.IP, port int) Peer {
	return Peer{
		Name:    strings.Replace(name, ".local.", "", 1),
		Address: address,
		Port:    port,
	}
}
