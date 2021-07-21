package file

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
)

// transferChunkSize is the number of bytes send on each write/read operation
// during the file transfer.
const transferChunkSize = 65536 //64kb

// OnProgressChange represents a callback to be executed when the file.Stream
// reports progress.
type OnProgressChange func(transferred int)

// Stream will copy the in content to the out.
//
// If the provided context has cancelation, dealine or timeout the stream
// will be interrupted.
//
// The optional argument onProg is callback function that is executed everytime
// a new chuck of data is transferred with success. If not needed nil can be sent.
//
// At the end it will return the total of bytes transferred and an error if any.
// If there is an error, it can be because the context got interrupted, an
// error reading the input content or writing to the output.
func Stream(ctx context.Context, in io.Reader, out io.Writer, onProg OnProgressChange) (int, error) {
	var transferred int
	buf := make([]byte, transferChunkSize)
	for {
		if ctx.Err() != nil {
			return -1, fmt.Errorf("file stream interrupted: %v", ctx.Err())
		}

		rc, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return -1, fmt.Errorf("file stream error read file: %v", err)
		}

		if rc == 0 {
			break
		}

		wc, err := out.Write(buf[:rc])
		if err != nil {
			return -1, fmt.Errorf("file stream error write file: %v", err)
		}
		transferred += wc
		if onProg != nil {
			onProg(transferred)
		}
	}
	return transferred, nil
}

// Lookup will get the file information based on the provided full path.
//
// Returns the file name and size.
//
// If there is an error, it can be because the file doesn't exists not possible
// to access or if the path points to a folder or a non regular file.
func Lookup(fileFullPath string) (string, int64, error) {
	cleanPath := filepath.Clean(fileFullPath)

	st, err := os.Stat(cleanPath)
	if err != nil {
		return "", -1, fmt.Errorf("file lookup error getting file info: %v", err)
	}

	if st.IsDir() || !st.Mode().IsRegular() {
		return "", -1, fmt.Errorf("file lookup error chekcing the file: file is not valid")
	}

	return filepath.Base(cleanPath), st.Size(), nil
}

// Checksum will make an SHA256 of the file content.
//
// It opens the file in read-only mode, used the file.Stream copy the content
// to the hash and return the string representation of the hash.
//
// If there is an error, it can be because there was an error opening the file
// or an error on streaming the file content to the hash.
func Checksum(ctx context.Context, fileFullPath string) (string, error) {
	f, err := os.OpenFile(filepath.Clean(fileFullPath), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("file checksum error opening file: %v", err)
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			clog.Error(cErr)
		}
	}()

	hash := sha256.New()

	if _, err = Stream(ctx, f, hash, nil); err != nil {
		return "", fmt.Errorf("file checksum error getting file content: %v", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
