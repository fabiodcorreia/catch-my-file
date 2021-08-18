package transfer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
	"github.com/fabiodcorreia/catch-my-file/pkg/file"
	"github.com/fabiodcorreia/catch-my-file/pkg/network"
	"github.com/fabiodcorreia/catch-my-file/pkg/transfer/internal/protocol"
)

const (
	// Time to wait to establish a connection with the receiver.
	dialTimeout = 10 //seconds
	// Time to wait to send the transfer request.
	writeTimeout = 30 //seconds
)

// ErrRejected signals that transfer was rejected.
var ErrRejected = errors.New(`REJECTED`)

// SendTransferReq receives a transfer, generates a request transfer
// message and send it to the receiver.
//
// If there is an error, it can be because it wasn't possible to establish
// a connection with the receiver, an error setting the timeout or an error
// when writing the message to the receiver.
func SendTransferReq(ctx context.Context, t *Transfer) (net.Conn, error) {

	rm := protocol.RequestMessage{
		FileName: t.FileName,
		FileSize: t.FileSize,
		Hostname: t.SenderName,
		Checksum: t.FileChecksum,
	}

	conn, err := net.DialTimeout(network.Type, t.SenderAddr.String(), dialTimeout*time.Second)
	if err != nil {
		return nil, fmt.Errorf("sender send transfer request error connecting: %v", err)
	}

	if err = conn.SetWriteDeadline(time.Now().Add(writeTimeout * time.Second)); err != nil {
		if err = conn.Close(); err != nil {
			return conn, err
		}
		return nil, fmt.Errorf("sender send transfer request set timeout error: %v", err)
	}

	if err = protocol.WriteRequestMessage(rm, conn); err != nil {
		if cErr := conn.Close(); cErr != nil {
			return conn, cErr
		}
		return nil, err
	}

	//Reset the write deadline set on SendTransferReq.
	if err = conn.SetWriteDeadline(time.Time{}); err != nil {
		return nil, fmt.Errorf("sender send transfer request write dealine error: %v", err)
	}

	return conn, nil
}

// WaitConfirmation will wait until the receiver accepts or rejects the transfer.
//
// If rejected it will just terminate and update the transfer status. Otherwise
// it will start sending the file content to the receiver.
func WaitConfirmation(ctx context.Context, i int, inOut io.ReadWriter, store *TransferStore) error {

	accept, err := protocol.ReadDecision(inOut)
	if err != nil {
		return fmt.Errorf("sender wait confirmation read decision error: %v", err)
	}

	trans := store.Get(i)

	if !accept {
		return ErrRejected
	}

	trans.Status = Accepted
	store.Update(i, trans)

	r, err := file.Open(trans.LocalFilePath, file.OPEN_READ)
	if err != nil {
		return fmt.Errorf("sender wait confirmation open file to send error: %v", err)
	}

	defer func() {
		if err = r.Close(); err != nil {
			clog.Error(fmt.Errorf("sender wait confirmation close file to send error: %v", err))
		}
	}()

	_, err = file.Stream(ctx, r, inOut, func(transferred int) {
		store.UpdateProgress(i, float64(transferred)/float64(trans.FileSize))
	})

	return err
}
