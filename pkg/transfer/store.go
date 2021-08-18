package transfer

import (
	"sync"
)

// OnStoreChange is a function that is executed everytime a
// transfer is added, removed or updated to the store.
type OnStoreChange func(i int)

// TransferStore is a thread-safe store that allows to store, retrieve, remove,
// and update transfer and also get notification when the content of the store changes.
type TransferStore struct {
	OnStoreChange
	mu   sync.Mutex
	data []*Transfer
}

// NewTransferStore will create a new instance of TransferStore which is thread-safe.
func NewStore() *TransferStore {
	return &TransferStore{
		data: make([]*Transfer, 0, 3),
	}
}

// Get will return an immutable transfer by the position i on the store.
//
// This transfer instance is a copy of the instance stored and for that
// reason any changes to the returned instance will not have effect on the
// stored instance.
func (s *TransferStore) Get(i int) *Transfer {
	s.mu.Lock()

	if i < 0 || i > s.Size()-1 {
		return nil
	}

	//Make a copy of the object stored to avoid outside mutations.
	cp := new(Transfer)
	*cp = *s.data[i]

	s.mu.Unlock()

	return cp
}

// Add will add a transfer to the store and return is position. If t is ni
// the transfer is not added and returns -1.
//
// Also executes the function OnStoreChange after the transfer gets added.
func (s *TransferStore) Add(t *Transfer) int {
	if t == nil {
		return -1
	}
	s.mu.Lock()
	s.data = append(s.data, t)
	i := len(s.data) - 1
	s.mu.Unlock()

	if s.OnStoreChange != nil {
		s.OnStoreChange(i)
	}

	return i
}

// Update will update the tranfer stored at position i with the values of t.
//
// Only the status, localfilepath and error are updated. It Also executes the
// function OnStoreChange after the transfer gets updated.
//
// If the status changes from Waiting it will close the waiting channel.
//
// If the status change to a final status it will close the progress channel.
func (s *TransferStore) Update(i int, t *Transfer) {
	if i < 0 || i > s.Size()-1 || t == nil {
		return
	}

	s.mu.Lock()
	s.update(i, t)
	s.mu.Unlock()

	if s.OnStoreChange != nil {
		s.OnStoreChange(i)
	}
}

func (s *TransferStore) update(i int, t *Transfer) {
	s.data[i].Status = t.Status
	s.data[i].LocalFilePath = t.LocalFilePath
	s.data[i].err = t.err

	// If the waiting channel is open and the status is not waiting,
	// it will close the channel to unlock the waiting.
	if s.data[i].wait != nil && t.Status != Waiting {
		close(s.data[i].wait)
		s.data[i].wait = nil
	}

	if s.data[i].Status.IsFinal() && s.data[i].prog != nil {
		close(s.data[i].prog)
		s.data[i].prog = nil
	}
}

// Size will return the current number of elements on the store.
func (s *TransferStore) Size() int {
	return len(s.data)
}

// AddToWait will add a transfer to the store, return is position and it will
// also return a channel that allows to wait until this transfer status changes
// from waiting to another status.
func (s *TransferStore) AddToWait(t *Transfer) (int, <-chan interface{}) {
	id := s.Add(t)
	if id == -1 {
		return id, nil
	}
	return id, s.data[id].waitDecision()
}

// FollowProgress returns a channel for the transfer on the specified position
// that allow the consumer to get the current progress of the transfer
// everytime it changes.
func (s *TransferStore) FollowProgress(i int) <-chan float64 {
	if i < 0 || i > s.Size()-1 {
		return nil
	}
	return s.data[i].progress()
}

// UpdateProgress allows to update the progress of the transfer on the
// specified index.
//
// Everythime this function is called the channel from FollowProgress
// will get the data provided here.
func (s *TransferStore) UpdateProgress(id int, progress float64) {
	s.data[id].progress() <- progress
}
