package transfer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileTransfer struct {
	FileFullPath string
	FileName     string
	FileSize     int64
	FileExt      string
	checksum     string
}

func NewFileTransfer(filePath string) (*FileTransfer, error) {
	st, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if st.IsDir() || !st.Mode().IsRegular() {
		return nil, fmt.Errorf("file is not valid")
	}

	return &FileTransfer{
		FileFullPath: filePath,
		FileName:     filepath.Base(filePath),
		FileSize:     st.Size(),
		FileExt:      filepath.Ext(filePath),
	}, nil
}

func (ft *FileTransfer) Open() (io.ReadCloser, error) {
	return os.OpenFile(ft.FileFullPath, os.O_RDONLY, os.ModePerm)
}

func (ft *FileTransfer) Checksum() (string, error) {
	if ft.checksum != "" {
		return ft.checksum, nil
	}

	f, err := ft.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	chks, err := fileChecksum(f)
	ft.checksum = chks
	return chks, err
}
