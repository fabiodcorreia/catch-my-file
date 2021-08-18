package peer

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
	"github.com/fabiodcorreia/catch-my-file/pkg/network"
	"github.com/grandcat/zeroconf"
)

const (
	// Zeroconf service name.
	serviceName = `_catchmyfile._tcp`
	// Zeroconf service domain we use local network.
	serviceDomain = `local.`
)

// Done is channel used to singal the termination of the service.
type Done chan<- interface{}

// PeerServer defines the server that registers the peer and discover other peers.
// It wrappes the logic for zeroconf.
type PeerServer struct {
	name  string
	port  int
	store *PeerStore
}

// NewServer will create a new peer server instace.
//
// The name to be added to the instance name, the port is the port number
// used for the TCP connections from where the files will be transferred.
func NewServer(name string, port int, store *PeerStore) *PeerServer {
	return &PeerServer{
		name:  name,
		port:  port,
		store: store,
	}
}

// Run will start the Server that will register the peer and start looking for
// peers on the local network with zeroconf.
//
// Each peer that is discovered will be added to the peer store.
//
// If there is an error, it can be because the server couldn't make the
// registrations, the server listener couldn't start or fail to start
// the discovery process.
func (s *PeerServer) Run(ctx context.Context, done Done) error {
	instance := fmt.Sprintf("catch-%s", s.name)
	sv, err := zeroconf.Register(instance, serviceName, serviceDomain, s.port, nil, nil)
	if err != nil {
		return fmt.Errorf("peer server run error to register: %v", err)
	}

	resolver, err := zeroconf.NewResolver(zeroconf.SelectIPTraffic(zeroconf.IPv4))
	if err != nil {
		return fmt.Errorf("peer server run error on start listening: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	if err = resolver.Browse(ctx, serviceName, serviceDomain, entries); err != nil {
		close(entries)
		return fmt.Errorf("peer server run error to discover: %v", err)
	}

	go func(results <-chan *zeroconf.ServiceEntry) {
		convEntry(results, s.store)
		sv.Shutdown()
		close(done)
		clog.Info("Peer server is closed")
	}(entries)

	return nil
}

// convEntry will grab each entry received from results channel, convert it
// into a Peer instance and add it to the store.
func convEntry(results <-chan *zeroconf.ServiceEntry, store *PeerStore) {
	if results != nil {
		for entry := range results {
			name := entry.HostName[:strings.Index(entry.HostName, `.`)]

			if len(entry.AddrIPv4) == 0 {
				clog.Error(fmt.Errorf("peer server conving entry error finding ipv4 for peer:%v", entry))
				continue
			}

			ipAddr := entry.AddrIPv4[0]
			addr, err := net.ResolveTCPAddr(network.Type, fmt.Sprintf("%v:%d", ipAddr, entry.Port))
			if err != nil {
				clog.Error(fmt.Errorf("peer server conving entry error resolving peer address:%v", err))
				continue
			}
			store.Add(newPeer(name, ipAddr, entry.Port, addr))
		}
	}
}
