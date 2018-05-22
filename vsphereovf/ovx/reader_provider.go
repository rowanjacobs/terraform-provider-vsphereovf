package ovx

import (
	"fmt"
	"io"
	"path/filepath"
)

//go:generate counterfeiter . ReaderProvider
type ReaderProvider interface {
	Reader(string) (io.Reader, int64, error)
	io.Closer
}

func NewReaderProvider(path string) (ReaderProvider, error) {
	if filepath.Ext(path) == ".ova" {
		return NewOVAReaderProvider(path)
	} else if filepath.Ext(path) == ".ovf" {
		return NewOVFReaderProvider(path)
	}
	return nil, fmt.Errorf("file '%s' does not have .ova or .ovf extension", filepath.Base(path))
}
