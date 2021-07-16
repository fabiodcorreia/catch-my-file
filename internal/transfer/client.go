package transfer

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/fabiodcorreia/catch-my-file/internal/store"
)

type Client struct {
	ServerAddress net.IP
	Port          int
	ft            *FileTransfer
	conn          net.Conn
}

func NewClient(addr net.IP, port int, ft *FileTransfer) *Client {
	return &Client{
		ServerAddress: addr,
		Port:          port,
		ft:            ft,
	}
}

func (tc *Client) SendRequest() error {
	ck, err := tc.ft.Checksum()
	if err != nil {
		return err
	}

	ckm, err := fillMessage(ck, transferChecksumBufferLen)
	if err != nil {
		return err
	}

	sz, err := fillMessage(strconv.FormatInt(tc.ft.FileSize, 10), transferSizeBufferLen)
	if err != nil {
		return err
	}

	fm, err := fillMessage(tc.ft.FileName, transferNameLen)
	if err != nil {
		return err
	}

	h, err := os.Hostname()
	if err != nil {
		return err
	}

	hn, err := fillMessage(h, transferHostnameLen)
	if err != nil {
		return err
	}

	message := make([]byte, 0, messageSize())
	message = append(message, append(ckm, append(sz, append(fm, hn...)...)...)...)

	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", tc.ServerAddress, tc.Port), time.Second*10)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	tc.conn = conn

	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	_, err = conn.Write(message)
	if err != nil {
		tc.conn.Close()
		return err
	}

	return nil
}

func (tc *Client) WaitSendOrStop(trans *store.Transfer) {
	defer tc.conn.Close()

	bufferConfirm := make([]byte, 1)
	_, err := tc.conn.Read(bufferConfirm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if string(bufferConfirm) != "1" {
		fmt.Println("file rejected")
		return
	}

	r, err := tc.ft.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer r.Close()

	sendBuffer := make([]byte, transferBufferLen)

	for {
		rc, err := r.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read error:" + err.Error())
			return
		}
		tc.conn.SetWriteDeadline(time.Time{})
		w, err := tc.conn.Write(sendBuffer[:rc])
		if err != nil {
			fmt.Println("client:" + err.Error())
			return
		}
		trans.UpdateTransferred(float64(w))
	}
}
