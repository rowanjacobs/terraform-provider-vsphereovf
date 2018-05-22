package ovx

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
)

type OVAReaderProvider struct {
	tarReader *tar.Reader
}

func NewOVAReaderProvider(path string) (OVAReaderProvider, error) {
	f, err := os.Open(path)
	if err != nil {
		return OVAReaderProvider{}, fmt.Errorf("OVA file '%s' not found", path)
	}

	return OVAReaderProvider{
		tarReader: tar.NewReader(f),
	}, nil
}

func (o OVAReaderProvider) Reader(relativePath string) (io.Reader, int64, error) {
	for {
		header, err := o.tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, fmt.Errorf("error reading OVA file: %s", err)
		}

		matched, err := path.Match(relativePath, path.Base(header.Name))
		if err != nil {
			return nil, 0, fmt.Errorf("error reading OVA file: %s", err)
		}

		if matched {
			return o.tarReader, header.Size, nil
		}
	}

	return nil, 0, fmt.Errorf("path '%s' not found in provided OVA", relativePath)
}

func (o OVAReaderProvider) Close() error {
	return nil
}
