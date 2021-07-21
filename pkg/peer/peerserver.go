package peer

import (
	"context"
	"fmt"
	"strings"

	"github.com/grandcat/zeroconf"
)

// Server defines the server that registers the peer and discover other peers.
// It wrappes the logic for zeroconf.
//
// Done is a channel that is closed when the Server finish the execution.
type Server struct {
	Done chan interface{}
	name string
	port int
}

// NewServer will create a new peer discover server.
//
// The port is the port number used for the TCP connections
// from where the files will be transferred.
func NewServer(name string, port int) *Server {
	return &Server{
		name: name,
		port: port,
		Done: make(chan interface{}),
	}
}

// Run will start the Server that will register the peer and start looking for
// peers on the local network with zeroconf.
//
// Arguments:
//
// - ctx is a context with cancel, what will be used to terminate the server.
//
// - peers is a channel of Peer that will stream each peer discovered on the
// network.
//
// Errors:
//
// - on registering the service
//
// - on start listening
//
// - on discover
func (s *Server) Run(ctx context.Context, peers chan<- Peer) error {
	instance := fmt.Sprintf("catch-%s", s.name)
	sv, err := zeroconf.Register(instance, serviceName, serviceDomain, s.port, nil, nil)
	if err != nil {
		return fmt.Errorf("peer server fail to register: %v", err)
	}

	resolver, err := zeroconf.NewResolver(zeroconf.SelectIPTraffic(zeroconf.IPv4))
	if err != nil {
		return fmt.Errorf("peer server fail to start listening: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	if err := resolver.Browse(ctx, serviceName, serviceDomain, entries); err != nil {
		close(entries)
		return fmt.Errorf("peer server fail to discover: %v", err)
	}

	go func(results <-chan *zeroconf.ServiceEntry) {
		convEntry(results, peers)
		sv.Shutdown()
		close(s.Done)
		close(peers)
		// entries is already closed by the context cancelation
	}(entries)

	return nil
}

// convEntry will grab each entry received from results channel, convert it
// into a Peer struct instance and send it to the peers channel.
func convEntry(results <-chan *zeroconf.ServiceEntry, peers chan<- Peer) {
	for entry := range results {
		name := strings.Replace(entry.HostName, `.local.`, ``, 1)
		ipAddr := entry.AddrIPv4[0] //TODO select the preferred IP address
		peers <- newPeer(name, ipAddr, entry.Port)
	}
}
