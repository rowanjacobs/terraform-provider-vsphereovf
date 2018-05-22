package ovx

import (
	"io"
	"os"
	"path/filepath"
)

type OVFReaderProvider struct {
	ovfDir string
}

func NewOVFReaderProvider(path string) (OVFReaderProvider, error) {
	_, err := os.Stat(path)
	if err != nil {
		return OVFReaderProvider{}, err
	}

	return OVFReaderProvider{
		ovfDir: filepath.Dir(path),
	}, nil
}

func (o OVFReaderProvider) Reader(relativePath string) (io.Reader, int64, error) {
	f, err := os.Open(filepath.Join(o.ovfDir, relativePath))
	if err != nil {
		return nil, 0, err
	}

	fInfo, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	return f, fInfo.Size(), nil
}

func (o OVFReaderProvider) Close() error {
	return nil
}
