package move

import "io"

type Src interface {
	io.Reader
	PrepareToRead(offset int64) error
	CompleteRead() error
}

type Dst interface {
	io.Writer
	PrepareToWrite() (offset int64, err error)
	CompleteWrite() error
}
