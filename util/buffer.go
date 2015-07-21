package util

import (
	"io"
)

type Buffer struct {
	io.Reader
	io.Writer
	io.Closer
}
