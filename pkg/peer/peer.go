package peer

import (
	"net"
)

// Peer defines a network peer that can send and receive files.
type Peer struct {
	// Name is the peers name.
	Name string
	// IPAddress is the net.IP address of the peer.
	IPAddress net.IP
	// Port is the network port where the peer will receive connections.
	Port int
	// Address is the resolved TCP Address of the peer IP+Port
	Address net.Addr
}

// newPeer will create a new instance of Peer struct and return it.
func newPeer(name string, ipAddress net.IP, port int, addr net.Addr) *Peer {
	return &Peer{
		Name:      name,
		IPAddress: ipAddress,
		Port:      port,
		Address:   addr,
	}
}
