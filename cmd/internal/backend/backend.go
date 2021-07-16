package backend

import (
	"context"
	"fmt"

	"github.com/fabiodcorreia/catch-my-file/internal/discover"
	"github.com/fabiodcorreia/catch-my-file/internal/store"
	"github.com/fabiodcorreia/catch-my-file/internal/transfer"
	"github.com/grandcat/zeroconf"
)

type appState uint

const (
	// All the components are off
	down appState = iota
	// Broadcast Server is running
	discoverServerUp
	// Broadcast Client is running
	discoverClientUp
	// Transfer Server is running
	transferServerUp
)

type Engine struct {
	ds     *discover.Server
	dc     *discover.Client
	ts     *transfer.Server
	ctx    context.Context
	cancel context.CancelFunc
	state  appState
	stream chan store.Peer
}

func NewEngine(hostname string, transferPort int) *Engine {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	e := &Engine{
		ds:     discover.NewServer(ctx, hostname, transferPort),
		dc:     discover.NewClient(ctx),
		ts:     transfer.NewServer(ctx, transferPort),
		state:  down,
		ctx:    ctx,
		cancel: cancel,
		stream: make(chan store.Peer, 1),
	}

	return e
}

func (e *Engine) DiscoverPeers() chan store.Peer {
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			e.stream <- store.NewPeer(entry.HostName, entry.AddrIPv4[0], entry.Port)
		}
	}(e.dc.PeerStream)
	return e.stream
}

func (e *Engine) ReceiveTransferNotification() chan *store.Transfer {
	return e.ts.TransferStream
}

func (e *Engine) DiscoverClientError() chan error {
	return e.ds.Err
}

func (e *Engine) DiscoverServerError() chan error {
	return e.dc.Err
}

func (e *Engine) TransferServerError() chan error {
	return e.ts.Err
}

func (e *Engine) Start() error {
	err := e.bootstrap()

	if err != nil {
		e.Shutdown()
		return err
	}

	return err
}

// Shutdown will cancel the context to finish all the backend components.
func (e *Engine) Shutdown() {
	e.cancel()
	// Wait until brocast server exits
	if e.state >= discoverServerUp {
		<-e.ds.Exit
	}
	// Wait until brocast client exits
	if e.state >= discoverClientUp {
		<-e.dc.Exit
	}
	// Wait until transfer exits
	if e.state >= transferServerUp {
		e.ts.Stop()
		<-e.ts.Exit
	}
}

// bootstrap will start a goroutine for each component and wait for the startup to finish.
// For each component ready it will also update the state of the application.
//
// If a component fails it will interrupt the sequence and return an error.
func (e *Engine) bootstrap() error {
	go e.ds.Run()

	err := waitForReady(e.ds.Ready, e.ds.Err)
	if err != nil {
		return fmt.Errorf("discover server fail before ready: %w", err)
	}

	go e.dc.Run()

	err = waitForReady(e.dc.Ready, e.dc.Err)
	if err != nil {
		return fmt.Errorf("discover client fail before ready: %w", err)
	}

	go e.ts.Run()

	err = waitForReady(e.ts.Ready, e.ts.Err)
	if err != nil {
		return fmt.Errorf("transfer client fail before ready: %w", err)
	}

	return nil
}

func waitForReady(ready chan interface{}, err chan error) error {
	select {
	case <-ready:
		return nil
	case e := <-err:
		return e
	}
}
