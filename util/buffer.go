package util

import (
  "io"
)

type buffer struct {
  io.Reader
  io.Writer
  io.Closer
}

