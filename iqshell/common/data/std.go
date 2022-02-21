package data

import (
	"io"
	"os"
)

var sdtout io.WriteCloser = os.Stdout
var sdterr io.WriteCloser = os.Stderr

func Stdout() io.WriteCloser {
	return sdtout
}

func Stderr() io.WriteCloser {
	return sdterr
}

func SetStdout(o io.WriteCloser) {
	sdtout = o
}

func SetStderr(e io.WriteCloser) {
	sdterr = e
}
