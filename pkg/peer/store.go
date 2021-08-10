package peer

import (
	"sync"
)

// OnPeerStoreChange is a function that is executed everytime a
// peer is added to the store.
type OnPeerStoreChange func(i int)

// PeerStore is a thread-safe store that allows to store and retrieve
// peers and also get notification when the content of the store changes.
type PeerStore struct {
	OnPeerStoreChange
	mu   sync.Mutex
	data []*Peer
}

// NewStore create a new instance of PeerStore.
func NewStore() *PeerStore {
	return &PeerStore{
		data: make([]*Peer, 0, 3),
	}
}

// Get return a peer from the specified index.
//
// If the index doesn't exists returns nil.
func (s *PeerStore) Get(i int) *Peer {
	s.mu.Lock()
	defer s.mu.Unlock()

	if i >= len(s.data) {
		return nil
	}

	return s.data[i]
}

// Add will append a peer to the existing list of peers.
//
// It returns the index where the peer whas stored.
//
// After the peer gets added the function OnPeerStoreChanged is executed.
func (s *PeerStore) Add(p *Peer) int {
	i := s.add(p)

	if s.OnPeerStoreChange != nil {
		s.OnPeerStoreChange(i)
	}

	return i
}

// Size returns the length of the store.

func (s *PeerStore) Size() int {
	return len(s.data)
}

// add will append a peer to the store using a mutext safe guard.
//
// Returns the index where the peer was stored.
func (s *PeerStore) add(p *Peer) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, p)
	return len(s.data) - 1

}
