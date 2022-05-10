package data

import (
	"io"
	"os"
)

var stdout io.WriteCloser = os.Stdout
var stderr io.WriteCloser = os.Stderr

func Stdout() io.WriteCloser {
	return stdout
}

func Stderr() io.WriteCloser {
	return stderr
}

func SetStdout(o io.WriteCloser) {
	stdout = o
}

func SetStderr(e io.WriteCloser) {
	stderr = e
}
