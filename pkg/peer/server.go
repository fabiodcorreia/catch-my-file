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
type PeerServer struct {
	name  string
	port  int
	store *PeerStore
}

// NewServer will create a new peer discover server.
//
// The port is the port number used for the TCP connections
// from where the files will be transferred.
func newServer(name string, port int, store *PeerStore) *PeerServer {
	return &PeerServer{
		name:  name,
		port:  port,
		store: store,
	}
}

// Run will start the Server that will register the peer and start looking for
// peers on the local network with zeroconf.
//
// Each peer that is dicovered will be added to the peer store.
func (s *PeerServer) Run(ctx context.Context, done chan<- interface{}) error {
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
		convEntry(results, s.store)
		sv.Shutdown()
		close(done)
	}(entries)

	return nil
}

// convEntry will grab each entry received from results channel, convert it
// into a Peer instance and add it to the store.
func convEntry(results <-chan *zeroconf.ServiceEntry, store *PeerStore) {
	for entry := range results {
		name := strings.Replace(entry.HostName, `.local.`, ``, 1)
		ipAddr := entry.AddrIPv4[0] //TODO select the preferred IP address
		store.Add(newPeer(name, ipAddr, entry.Port))
	}
}
