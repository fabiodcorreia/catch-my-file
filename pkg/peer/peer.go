package peer

import (
	"net"
)

// Peer defines a network peer that can send and receive files.
type Peer struct {
	Name      string   // Name is the peers name.
	IPAddress net.IP   // IPAddress is the net.IP address of the peer.
	Port      int      // Port is the network port where the peer will receive connections.
	Address   net.Addr // Address is the resolved TCP Address of the peer IP+Port.
	Me        bool     // Me identify the peer as the local peer.
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
