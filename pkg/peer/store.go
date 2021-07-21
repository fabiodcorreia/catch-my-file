package peer

import (
	"sync"
)

type OnPeerStoreChange func(i int)

type PeerStore struct {
	mu      sync.Mutex
	data    []*Peer
	actions []OnPeerStoreChange
}

// NewPeerStore will create a new instance of PeerStore which is thread-safe.
func newStore() *PeerStore {
	return &PeerStore{
		data:    make([]*Peer, 0, 3),
		actions: make([]OnPeerStoreChange, 0),
	}
}

// Get will return a peer by the position i on the store.
func (s *PeerStore) Get(i int) *Peer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[i]
}

// Add will add a peer to the store.
func (s *PeerStore) Add(p *Peer) {
	s.mu.Lock()
	s.data = append(s.data, p)

	s.mu.Unlock()
	s.fireChange(len(s.data) - 1)

}

// Remove will remove a peer by position i on the store.
func (s *PeerStore) Remove(i int) {
	s.mu.Lock()
	copy(s.data[i:], s.data[i+1:])
	s.data[s.Size()-1] = nil
	s.data = s.data[:len(s.data)-1]

	s.mu.Unlock()
	s.fireChange(i)
}

// Size will return the current number of elements on the store.
func (s *PeerStore) Size() int {
	return len(s.data)
}

// addOnChangeListener will add a new OnPeerStoreChange function that
// will be executed everytime the store content changes.
func (s *PeerStore) AddOnChangeListener(action OnPeerStoreChange) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.actions = append(s.actions, action)
}

// fireChange will execute each OnPeerStoreChange function registred
func (s *PeerStore) fireChange(i int) {
	if len(s.actions) > 0 {
		for _, af := range s.actions {
			af(i)
		}
	}
}
