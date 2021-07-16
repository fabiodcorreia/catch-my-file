package discover

import (
	"context"

	"github.com/grandcat/zeroconf"
)

type Client struct {
	PeerStream chan *zeroconf.ServiceEntry
	Ready      chan interface{}
	Exit       chan interface{}
	Err        chan error
	ctx        context.Context
}

func NewClient(ctx context.Context) *Client {
	return &Client{
		PeerStream: make(chan *zeroconf.ServiceEntry),
		Ready:      make(chan interface{}),
		Exit:       make(chan interface{}, 1),
		Err:        make(chan error, 1),
		ctx:        ctx,
	}
}

func (bc *Client) stop(err error) {
	bc.Err <- err
	close(bc.Err)
	close(bc.Exit)
}

func (bc *Client) Run() {

	resolver, err := zeroconf.NewResolver(zeroconf.SelectIPTraffic(zeroconf.IPv4))
	if err != nil {
		bc.stop(err)
		return
	}

	zeroconf.SelectIPTraffic(zeroconf.IPv4)

	close(bc.Ready)

	err = resolver.Browse(bc.ctx, serviceName, serviceDomain, bc.PeerStream)
	if err != nil {
		bc.Err <- err
	}

	<-bc.ctx.Done()
	bc.stop(err)
}
