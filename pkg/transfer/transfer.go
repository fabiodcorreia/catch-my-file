package transfer

import (
	"net"
)

// Status represents the Transfer status.
type Status int

const (
	// Transfer is waiting for confirmation.
	Waiting Status = iota + 1
	// Transfer was accepted.
	Accepted
	// Transfer is getting verified.
	Verifying
	// Transfer was rejected.
	Rejected
	// Transfer is completed.
	Completed
	// Transfer is error state.
	Error
)

// IsFinal returns true if the status is a final status, which
// means it will not change anymore.
func (s Status) IsFinal() bool {
	return (s == Error || s == Rejected || s == Completed)
}

// String convert a transfer status into a string representation
// of that status.
func (s Status) String() string {
	switch s {
	case Waiting:
		return `Waiting`
	case Verifying:
		return `Verifying`
	case Accepted:
		return `Accepted`
	case Rejected:
		return `Rejected`
	case Completed:
		return `Completed`
	case Error:
		return `Error`
	}
	return ``
}

// Direction represents the direction of the transfer.
//
// Upload or Download.
type Direction int

const (
	// Transfer is an Upload
	Upload Direction = iota + 1
	// Transfer is a Download
	Download
)

// Transfer wraps the transfer information
type Transfer struct {
	Direction
	Status        Status
	SenderName    string
	SenderAddr    net.Addr
	FileName      string
	FileChecksum  string
	FileSize      int64
	LocalFilePath string           // Full path to local file system. Sender/read path, Receiver/save path.
	wait          chan interface{} // Waiting channel used to notify when the user accepted or rejected.
	prog          chan float64     // Progress channel used to report the progress updates.
	err           error            // Error that occurred to the transfer.
}

// NewTransfer creates a new Transfer instance.
//
// It requires the name of the file, the checksum, the peers name
// that is sending or receiving depending if it's an upload or download,
// the size of the file, the address of the peer and the direction.
func NewTransfer(name, check, sender string, size int64, addr net.Addr, dir Direction) *Transfer {
	return &Transfer{
		Status:       Waiting,
		Direction:    dir,
		SenderName:   sender,
		SenderAddr:   addr,
		FileName:     name,
		FileChecksum: check,
		FileSize:     size,
	}
}

// SetError register an error to the transfer and changes the status to Error.
func (t *Transfer) SetError(err error) {
	t.Status = Error
	t.err = err
}

// Error return the last error registered on the transfer.
func (t *Transfer) Error() error {
	return t.err
}

// waitDecision will create and return a channel that will be used to signal
// when the transfer status changed from Waiting to anotehr status.
func (t *Transfer) waitDecision() <-chan interface{} {
	if t.wait == nil {
		t.wait = make(chan interface{})
	}
	return t.wait
}

// progress will create and return a channel that will be used to stream the
// amount of bytes already transferred.
func (t *Transfer) progress() chan float64 {
	if t.prog == nil {
		t.prog = make(chan float64)
	}
	return t.prog
}
