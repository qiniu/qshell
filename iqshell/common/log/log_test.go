package log

import (
	"testing"
)

func TestProgress(t *testing.T) {
	LoadConsole(LevelDebug)
	Debug("progress:%v", " debug")
}
