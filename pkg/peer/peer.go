package peer

import (
	"context"
	"net"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
func newPeer(name string, address net.IP, port int) *Peer {
	return &Peer{
		Name:    name,
		Address: address,
		Port:    port,
	}
}

type PeerController struct {
	store  *PeerStore
	server *PeerServer
	view   *PeerList
}

func New(hostname string, port int) *PeerController {
	s := newStore()
	return &PeerController{
		store:  s,
		server: newServer(hostname, port, s),
		view:   newListView(s),
	}
}

func (c *PeerController) View() *container.TabItem {
	return container.NewTabItemWithIcon("Peers", theme.ComputerIcon(), c.view)
}

func (c *PeerController) Start(ctx context.Context, done chan<- interface{}) error {
	return c.server.Run(ctx, done)
}
