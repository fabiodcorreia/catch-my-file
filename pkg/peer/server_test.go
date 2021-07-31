package peer

import (
	"net"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
)

func Test_convEntry(t *testing.T) {
	t.Run("entries is nil", func(t *testing.T) {

		convEntry(nil, nil)
	})

	t.Run("entries send one peer and it get updated on the store", func(t *testing.T) {
		entries := make(chan *zeroconf.ServiceEntry)
		store := NewStore()

		go convEntry(entries, store)

		entry := zeroconf.NewServiceEntry("instance", "service", "domain")
		entry.HostName = "peer-1.lan"
		entry.AddrIPv4 = append(entry.AddrIPv4, net.ParseIP("192.168.1.1"))
		entry.Port = 8822

		entries <- entry

		time.Sleep(500 * time.Millisecond)
		close(entries)

		if store.Get(0).Name != "peer-1" {
			t.Errorf("convEntry expected name = %v but got %v", "peer-1", store.Get(0).Name)
		}

		if store.Get(0).IPAddress.String() != "192.168.1.1" {
			t.Errorf("convEntry expected ip address = %v but got %v", "192.168.1.1", store.Get(0).Address.String())
		}

		if store.Get(0).Port != 8822 {
			t.Errorf("convEntry expected port = %v but got %v", "192.168.1.1", store.Get(0).Port)
		}

		if store.Get(0).Address.String() != "192.168.1.1:8822" {
			t.Errorf("convEntry expected port = %v but got %v", "192.168.1.1:8822", store.Get(0).Port)
		}
	})

	t.Run("entries send one peer with no IPv4", func(t *testing.T) {
		entries := make(chan *zeroconf.ServiceEntry)
		store := NewStore()

		go convEntry(entries, store)

		entry := zeroconf.NewServiceEntry("instance", "service", "domain")
		entry.HostName = "peer-1.lan"

		entries <- entry

		time.Sleep(500 * time.Millisecond)
		close(entries)

		if store.Size() != 0 {
			t.Errorf("convEntry store size expected = %v but got %v", 0, store.Size())
		}
	})

	t.Run("entries send one peer with invalid port", func(t *testing.T) {
		entries := make(chan *zeroconf.ServiceEntry)
		store := NewStore()

		go convEntry(entries, store)

		entry := zeroconf.NewServiceEntry("instance", "service", "domain")
		entry.HostName = "peer-1.lan"
		entry.AddrIPv4 = append(entry.AddrIPv4, net.ParseIP("192.168.1.1"))
		entry.Port = -1

		entries <- entry

		time.Sleep(500 * time.Millisecond)
		close(entries)

		if store.Size() != 0 {
			t.Errorf("convEntry store size expected = %v but got %v", 0, store.Size())
		}
	})

}
