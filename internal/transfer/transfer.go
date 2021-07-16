package transfer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/fabiodcorreia/catch-my-file/internal/store"
)

const (
	transferBufferLen         int = 65536 //64kb
	transferChecksumBufferLen int = 64
	transferSizeBufferLen     int = 10
	transferNameLen           int = 128
	transferHostnameLen       int = 32
)

func messageSize() int {
	return transferChecksumBufferLen + transferSizeBufferLen + transferNameLen + transferHostnameLen
}

func messageCheckIdx() int {
	return transferChecksumBufferLen
}

func messageSizeIdx() int {
	return messageCheckIdx() + transferSizeBufferLen
}

func messageNameIdx() int {
	return messageCheckIdx() + transferNameLen
}

func messageHostIdx() int {
	return messageNameIdx() + transferHostnameLen
}

func fileChecksum(r io.Reader) (string, error) {
	hash := sha256.New()

	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func WriteFile(ctx context.Context, r io.Reader, filePath string, t *store.Transfer) error {
	w, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil && err != io.EOF {
		return err
	}

	buf := make([]byte, transferBufferLen)
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		nr, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if nr == 0 {
			break
		}

		nw, err := w.Write(buf[:nr])
		if err != nil {
			return err
		}

		t.UpdateTransferred(float64(nw))
	}
	/*
		if t.Transferred() != t.Size {
			return fmt.Errorf("file size and transferred size are different")
		}*/
	return nil
}

func fillMessage(content string, toFill int) ([]byte, error) {
	if len(content) > toFill {
		return nil, fmt.Errorf("message is longer than the fill length")
	}

	buffer := make([]byte, toFill)
	copy(buffer, content)
	return buffer, nil
}

func trimMessage(message []byte) string {
	for i := range message {
		if message[i] == 0 {
			return string(message[:i])
		}
	}
	return string(message)
}
