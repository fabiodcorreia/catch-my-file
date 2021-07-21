package transfer

import (
	"fmt"
	"io"
	"strconv"
)

// The length of each message field in bytes.
const (
	fieldChecksumLen = 64
	fieldFileSizeLen = 10
	fieldFileNameLen = 128
	fieldHostnameLen = 32
)

// The end index of each message field, calculated using the field length
// plus the previous field end index.
const (
	idxFieldChecksum = fieldChecksumLen
	idxFieldFileSize = idxFieldChecksum + fieldFileSizeLen
	idxFieldFileName = idxFieldFileSize + fieldFileNameLen
	idxFieldHostname = idxFieldFileName + fieldHostnameLen
)

const (
	// messageRequestLen is to length of the full request transfer message.
	messageRequestLen = fieldChecksumLen + fieldFileSizeLen + fieldFileNameLen + fieldHostnameLen

	// messageTransferLen is the number of bytes send on each write/read operations
	// during the file transfer.
	messageTransferLen = 65536 //64kb
)

// requestMessage wraps the request message data that is sent and received  by
// the peers when a new file transfer is requested.
type requestMessage struct {
	FileName string
	FileSize int64
	Hostname string
	Checksum string
}

// writeRequestMessage will create a structured binary message to sent to
// the receiver storing it on the w.
//
// If any of the message fields is not valid nothing will be written.
func writeRequestMessage(m requestMessage, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("writing request message fail: writer is nil")
	}

	check := make([]byte, fieldChecksumLen)
	if err := fillMessageField(m.Checksum, check); err != nil {
		return fmt.Errorf("writing message request fail checksum field: %v", err)
	}

	size := make([]byte, fieldFileSizeLen)
	if err := fillMessageField(strconv.FormatInt(m.FileSize, 10), size); err != nil {
		return fmt.Errorf("writing message request fail size field: %v", err)
	}

	name := make([]byte, fieldFileNameLen)
	if err := fillMessageField(m.FileName, name); err != nil {
		return fmt.Errorf("writing message request fail file name field: %v", err)
	}

	host := make([]byte, fieldHostnameLen)
	if err := fillMessageField(m.Hostname, host); err != nil {
		return fmt.Errorf("writing message request fail hostname field: %v", err)
	}

	p := make([]byte, messageRequestLen)
	copy(p[:idxFieldChecksum], check)
	copy(p[idxFieldChecksum:idxFieldFileSize], size)
	copy(p[idxFieldFileSize:idxFieldFileName], name)
	copy(p[idxFieldFileName:idxFieldHostname], host)

	if _, err := w.Write(p); err != nil {
		return fmt.Errorf("writing message request error: %v", err)
	}

	return nil
}

// readRequestMessage will read a structured binary message and bind into
// a requestMessage instance.
//
// If any error occurs no data will be bound to the requestMessage instance.
func readRequestMessage(m *requestMessage, r io.Reader) error {
	bufferMessage := make([]byte, messageRequestLen)

	if r == nil {
		return fmt.Errorf("revert request message fail: reader is nil")
	}

	rc, err := r.Read(bufferMessage)
	switch {
	case err != nil:
		return fmt.Errorf("revert request message fail: %v", err)
	case rc != messageRequestLen:
		return fmt.Errorf("revert request message fail: message size not correct")
	}

	size, err := strconv.ParseInt(
		trimMessageField(bufferMessage[idxFieldChecksum:idxFieldFileName]),
		10,
		64,
	)
	if err != nil {
		return fmt.Errorf("revert request message fail to convert file size: %v", err)
	}

	m.FileSize = size
	m.FileName = trimMessageField(bufferMessage[idxFieldFileSize:idxFieldFileName])
	m.Hostname = trimMessageField(bufferMessage[idxFieldFileName:idxFieldHostname])
	m.Checksum = trimMessageField(bufferMessage[:idxFieldChecksum])

	return nil
}

// fillMessageField will receive a content string and convert it into a []byte
// filling the remaining positions of the []byte length with 0 value bytes.
//
// If the length of the content is larger than then length of the buffer an erro is returned.
func fillMessageField(content string, buffer []byte) error {
	if len(content) > len(buffer) {
		return fmt.Errorf("content is longer than field length: %d vs %d", len(content), len(buffer))
	}

	copy(buffer, content)
	return nil
}

// trimMessageField will look for the first 0 byte value on the field content
// and return a string with the field content before the 0 byte value.
//
// This will remove the extra 0 bytes used to fill the message field.
func trimMessageField(field []byte) string {
	for i := range field {
		if field[i] == 0 {
			return string(field[:i])
		}
	}
	return string(field)
}

// writeConfirmation will write one byte to the writer depending if the
// accept argument is true or false.
func writeConfirmation(accept bool, w io.Writer) error {
	if w == nil {
		return fmt.Errorf("write confirmation error: writeer is nil")
	}

	buffer := make([]byte, 1)
	switch accept {
	case true:
		buffer[0] = 1
	default:
		buffer[0] = 0
	}

	if _, err := w.Write(buffer); err != nil {
		return fmt.Errorf("write confirmation error: %v", err)
	}

	return nil
}

// readConfirmation will read one byte from the reader and return if
// the request was accepted or not.
func readConfirmation(r io.Reader) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("read confirmation error: reader is nil")
	}

	buffer := make([]byte, 1)
	if _, err := r.Read(buffer); err != nil {
		return false, fmt.Errorf("read confirmation error: %v", err)
	}

	if buffer[0] == 1 {
		return true, nil
	}
	return false, nil
}
