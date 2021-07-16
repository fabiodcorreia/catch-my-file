package store

import (
	"net"
	"strings"
	"time"
)

type Peer struct {
	Name      string
	Address   net.IP
	Port      int
	Timestamp time.Time
}

func NewPeer(name string, address net.IP, port int) Peer {
	return Peer{
		Name:      strings.Replace(name, ".local.", "", 1),
		Address:   address,
		Port:      port,
		Timestamp: time.Now(),
	}
}

type Transfer struct {
	Checksum         string
	Name             string
	Size             float64
	SourceAddr       net.Addr
	SourceName       string
	Accept           chan string
	IsToSend         bool
	transferred      float64
	progressCallback func(progress float64)
}

func NewTransfer(checksum string, name string, size float64, sourceName string, sourceAddr net.Addr) *Transfer {
	return &Transfer{
		Checksum:   checksum,
		Name:       name,
		Size:       size,
		SourceName: sourceName,
		SourceAddr: sourceAddr,
		Accept:     make(chan string),
	}
}

func (t *Transfer) OnProgressChange(callback func(progress float64)) {
	if callback != nil {
		t.progressCallback = callback
	}
}

func (t *Transfer) UpdateTransferred(value float64) {
	t.transferred += value
	if t.progressCallback != nil {
		t.progressCallback(t.transferred / t.Size)
	}
}

func (t *Transfer) Transferred() float64 {
	return t.transferred
}
