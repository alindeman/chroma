package chroma

import (
	"io"
)

type nopCloser struct {
	io.Reader
}

func (nc nopCloser) Close() error {
	return nil
}
