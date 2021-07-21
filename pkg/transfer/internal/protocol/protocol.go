package protocol

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

// messageRequestLen is to length of the full request transfer message.
const messageRequestLen = fieldChecksumLen + fieldFileSizeLen + fieldFileNameLen + fieldHostnameLen

// RequestMessage wraps the request message data that is sent and received  by
// the peers when a new file transfer is requested.
type RequestMessage struct {
	FileName string
	FileSize int64
	Hostname string
	Checksum string
}

// WriteRequestMessage will create a structured binary message to sent to
// the receiver writing it to the output.
//
// If any of the message fields is not valid nothing will be written.
//
// If there is an error, it can be because the writer was nil, the data
// provided for each field is not valid or an error writing to the output.
func WriteRequestMessage(m RequestMessage, out io.Writer) error {
	if out == nil {
		return fmt.Errorf("protocol write request message error: writer is nil")
	}

	check := make([]byte, fieldChecksumLen)
	if err := fillMessageField(m.Checksum, check); err != nil {
		return fmt.Errorf("protocol write message request error on field checksum: %v", err)
	}

	size := make([]byte, fieldFileSizeLen)
	if err := fillMessageField(strconv.FormatInt(m.FileSize, 10), size); err != nil {
		return fmt.Errorf("protocol write message request error on field size: %v", err)
	}

	name := make([]byte, fieldFileNameLen)
	if err := fillMessageField(m.FileName, name); err != nil {
		return fmt.Errorf("protocol write message request error on field name: %v", err)
	}

	host := make([]byte, fieldHostnameLen)
	if err := fillMessageField(m.Hostname, host); err != nil {
		return fmt.Errorf("protocol write message request error on field hostname: %v", err)
	}

	p := make([]byte, messageRequestLen)
	copy(p[:idxFieldChecksum], check)
	copy(p[idxFieldChecksum:idxFieldFileSize], size)
	copy(p[idxFieldFileSize:idxFieldFileName], name)
	copy(p[idxFieldFileName:idxFieldHostname], host)

	if _, err := out.Write(p); err != nil {
		return fmt.Errorf("protocol write message request error writing the output: %v", err)
	}

	return nil
}

// ReadRequestMessage will read a structured binary message and bind it into
// a requestMessage instance provided as argument.
//
// If there is an error, it can be because the reader was nil, an error
// reading the input, the length of the message is not correct or fail to
// convert file size to a valid number.
//
// If any error occurs no data will be bound to the requestMessage instance.
func ReadRequestMessage(m *RequestMessage, in io.Reader) error {
	bufferMessage := make([]byte, messageRequestLen)

	if in == nil {
		return fmt.Errorf("protocol read request message error: reader is nil")
	}

	rc, err := in.Read(bufferMessage)
	switch {
	case err != nil:
		return fmt.Errorf("protocol read request message error reading the input: %v", err)
	case rc != messageRequestLen:
		return fmt.Errorf("protocol read request message error: message size not correct")
	}

	size, err := strconv.ParseInt(
		trimMessageField(bufferMessage[idxFieldChecksum:idxFieldFileSize]),
		10,
		64,
	)
	if err != nil {
		return fmt.Errorf("read request message error converting file size: %v", err)
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
		return fmt.Errorf("only allowed %d characteres but found %d", len(buffer), len(content))
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

// WriteDecision will write one byte to the writer depending if the
// accept argument is true or false.
//
// If there is an error, it can be because the writer was nil or
// and error occurred writing to the output.
func WriteDecision(accept bool, out io.Writer) error {
	if out == nil {
		return fmt.Errorf("protocol write decision error: output is nil")
	}

	buffer := make([]byte, 1)
	switch accept {
	case true:
		buffer[0] = byte(1)
	default:
		buffer[0] = byte(0)
	}

	if _, err := out.Write(buffer); err != nil {
		return fmt.Errorf("protocol write decision writing to output: %v", err)
	}

	return nil
}

// ReadDecision will read one byte from the reader and return if
// the request was accepted or not.
//
// If there is an error, it can be because the reader was nil or
// and error occurred reading the input.
func ReadDecision(in io.Reader) (bool, error) {
	if in == nil {
		return false, fmt.Errorf("protocol read decision error: reader is nil")
	}

	buffer := make([]byte, 1)
	if _, err := in.Read(buffer); err != nil {
		return false, fmt.Errorf("protocol read decision error reading the input: %v", err)
	}

	if buffer[0] == 1 {
		return true, nil
	}
	return false, nil
}
