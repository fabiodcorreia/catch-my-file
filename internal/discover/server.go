package discover

import (
	"context"
	"fmt"

	"github.com/grandcat/zeroconf"
)

type Server struct {
	Name  string
	Port  int
	Ready chan interface{}
	Exit  chan interface{}
	Err   chan error
	ctx   context.Context
}

func NewServer(ctx context.Context, name string, port int) *Server {
	return &Server{
		Name:  name,
		Port:  port,
		Ready: make(chan interface{}),
		Exit:  make(chan interface{}, 1),
		Err:   make(chan error, 1),
		ctx:   ctx,
	}
}

func (bs *Server) Run() {
	server, err := zeroconf.Register(fmt.Sprintf("catch-%s", bs.Name), serviceName, serviceDomain, bs.Port, nil, nil)
	if err != nil {
		bs.Err <- err
		close(bs.Err)
		close(bs.Exit)
		return
	}

	defer server.Shutdown()

	close(bs.Ready)

	<-bs.ctx.Done()

	close(bs.Err)
	close(bs.Exit)
}
