package transfer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
	"github.com/fabiodcorreia/catch-my-file/pkg/file"
	"github.com/fabiodcorreia/catch-my-file/pkg/network"
	"github.com/fabiodcorreia/catch-my-file/pkg/transfer/internal/protocol"
)

// Done is channel used to singal the termination of the service.
type Done chan<- interface{}

type Receiver struct {
	port  int
	store *TransferStore
}

// NewReceiver will create a new Receiver server that will wait for
// connections on the specified port.
func NewReceiver(port int, store *TransferStore) *Receiver {
	return &Receiver{
		port:  port,
		store: store,
	}
}

// Run will start the receiver starting the listener to receive requests.
func (rv *Receiver) Run(ctx context.Context, done Done) error {
	listener, lErr := net.Listen(network.Type, fmt.Sprintf(":%d", rv.port))
	if lErr != nil {
		close(done)
		return fmt.Errorf("receiver run listen error: %v", lErr)
	}

	go waitForRequests(ctx, listener, done, rv.store)
	go watchdog(ctx, listener)

	return nil
}

// watchdog will wait until the context gets cancelled and after
// that it will stop the listener.
func watchdog(ctx context.Context, listener net.Listener) {
	<-ctx.Done()
	if err := listener.Close(); err != nil {
		clog.Error(err)
	}
}

// waitForRequests will wait for new connections from senders and for each
// connection will handle handle the request.
func waitForRequests(ctx context.Context, listener net.Listener, done Done, store *TransferStore) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			switch {
			case errors.Is(err, net.ErrClosed):
				clog.Info("Receiver server is closed")
			default:
				clog.Error(err)
			}
			close(done)
			return
		}
		go handleRequest(ctx, conn, store)
	}
}

// handleRequest will receive a new connection, add the transfer to the store,
// wait for the confirmation or rejection and if accepted start receiving and
// storing the file.
func handleRequest(ctx context.Context, conn net.Conn, store *TransferStore) {
	defer conn.Close()

	id, err := reqDecisionAndWait(ctx, store, conn, conn.RemoteAddr())
	if err != nil {
		clog.Error(err)
		return
	}

	trans := store.Get(id)

	clog.Info("writing decision to sender: %v", trans.Status)

	if err = protocol.WriteDecision(trans.Status == Accepted, conn); err != nil {
		trans.SetError(err)
		store.Update(id, trans)
		clog.Error(err)
		return
	}

	if trans.Status != Accepted {
		clog.Info("transfer rejected %d", id)
		return //It was rejected just end the work
	}

	// The stored file should get updated with the path where to store the file
	w, err := file.Open(trans.LocalFilePath, file.OPEN_WRITE)
	if err != nil {
		clog.Error(err)
		return
	}

	defer func() {
		if cErr := w.Close(); cErr != nil {
			clog.Error(cErr)
		}
	}()

	clog.Info("receiving file from sender and store at: %s", trans.LocalFilePath)

	rcvSize, err := file.Stream(ctx, conn, w, func(transferred int) {
		store.UpdateProgress(id, float64(transferred)/float64(trans.FileSize))
	})
	if err != nil {
		trans.SetError(err)
		store.Update(id, trans)
		clog.Error(err)
		return
	}

	trans.Status = Verifying
	store.Update(id, trans)

	verifyFile(ctx, trans, rcvSize, id, store)
}

func verifyFile(ctx context.Context, t *Transfer, rcvSize, id int, store *TransferStore) {
	f, err := file.Open(t.LocalFilePath, file.OPEN_READ)
	if err != nil {
		clog.Error(err)
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			clog.Error(err)
		}
	}()

	if err = verifyTransfer(ctx, t, rcvSize, f); err != nil {
		t.SetError(err)
	} else {
		t.Status = Completed
	}

	store.Update(id, t)
}

// reqDecisionAndWait will add the transfer to the store and wait for confirmation
// by the user or a cancelation of the context.
func reqDecisionAndWait(ctx context.Context, store *TransferStore, r io.Reader, addr net.Addr) (int, error) {
	var rm protocol.RequestMessage
	if err := protocol.ReadRequestMessage(&rm, r); err != nil {
		clog.Error(err)
		return -1, err
	}

	id, wait := store.AddToWait(NewTransfer(rm.FileName, rm.Checksum, rm.Hostname, rm.FileSize, addr, Download))

	clog.Info("waiting for trans: %d", id)

	select {
	case <-wait:
		return id, nil
	case <-ctx.Done():
		clog.Info("waiting interrupted: %d", id)
		return -1, ctx.Err()
	}
}

// verifyTransfer will verify is the amount of data transferred matches with
// the amount received and will check with the checkshum also match.
func verifyTransfer(ctx context.Context, trans *Transfer, rcvSize int, in io.Reader) error {
	if trans == nil {
		return fmt.Errorf("receiver verify transfer error: transfer is nil")
	}

	if int64(rcvSize) != trans.FileSize {
		return fmt.Errorf("receiver verify transfer error: data size don't match")
	}

	check, err := file.Checksum(ctx, in)
	if err != nil {
		return err
	}

	if check != trans.FileChecksum {
		return fmt.Errorf("receiver verify transfer error: checksum doesn't match")
	}
	return nil
}
