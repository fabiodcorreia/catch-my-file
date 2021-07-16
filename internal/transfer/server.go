package transfer

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/fabiodcorreia/catch-my-file/internal/store"
)

type Server struct {
	Port           int
	Ready          chan interface{}
	Exit           chan interface{}
	Err            chan error
	TransferStream chan *store.Transfer
	listener       net.Listener
	ctx            context.Context
}

func NewServer(ctx context.Context, port int) *Server {
	return &Server{
		Port:           port,
		Ready:          make(chan interface{}),
		Exit:           make(chan interface{}, 1),
		Err:            make(chan error, 1),
		TransferStream: make(chan *store.Transfer),
		ctx:            ctx,
	}

}

func (ts *Server) Stop() {
	ts.stop(nil)
}

func (ts *Server) stop(err error) {
	if err != nil {
		ts.Err <- err
	}
	err = ts.listener.Close()
	if err != nil {
		ts.Err <- err
	}
	close(ts.Err)
	close(ts.Exit)
}

func (ts *Server) Run() {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", ts.Port))
	if err != nil {
		ts.Err <- err
		return
	}

	ts.listener = listener

	close(ts.Ready)

	for {
		c, err := ts.listener.Accept()

		if err != nil {
			ts.Err <- err
			return
		}

		if ts.ctx.Err() != nil {
			ts.stop(ts.ctx.Err())
			return
		}

		go ts.transfer(c)
	}

}

func (ts *Server) transfer(conn net.Conn) {

	tName, tCheck, tHostname, tSize, err := receiveTransferData(conn)
	if err != nil {
		ts.Err <- err
		return
	}

	tt := store.NewTransfer(tCheck, tName, float64(tSize), tHostname, conn.RemoteAddr())

	ts.TransferStream <- tt

	filePath, open := <-tt.Accept

	if !open {
		ts.Err <- fmt.Errorf("transfer %s rejected", tt.Name)
		conn.Write([]byte("0"))
		return
	}

	conn.Write([]byte("1"))

	err = WriteFile(ts.ctx, conn, filePath, tt)
	if err != nil {
		ts.Err <- err
		return
	}

	err = verifyFile(filePath, tt.Checksum)
	if err != nil {
		ts.Err <- err
		return
	}
	ts.Err <- fmt.Errorf("file %s transferred and verified", tName)
}

func verifyFile(filePath string, checksum string) error {
	ft, err := NewFileTransfer(filePath)
	if err != nil {
		return err
	}

	c, err := ft.Checksum()

	if err != nil {
		return err
	}

	if c != checksum {
		err = fmt.Errorf("file checksum don't match")
		e := os.Remove(ft.FileFullPath)
		if e != nil {
			return fmt.Errorf("%s - %w", err, e)
		}
		return err
	}
	return nil
}

func receiveTransferData(conn net.Conn) (string, string, string, int64, error) {
	bufferMessage := make([]byte, messageSize())
	_, err := conn.Read(bufferMessage)
	if err != nil {
		return "", "", "", 0, err
	}

	bufferChecksum := bufferMessage[:messageCheckIdx()]
	bufferFileSize := bufferMessage[messageCheckIdx():messageSizeIdx()]
	bufferFileName := bufferMessage[messageSizeIdx():messageNameIdx()]
	bufferHostName := bufferMessage[messageNameIdx():messageHostIdx()]

	tSize, err := strconv.ParseInt(trimMessage(bufferFileSize), 10, 64)

	return trimMessage(bufferFileName), trimMessage(bufferChecksum), trimMessage(bufferHostName), tSize, err
}
