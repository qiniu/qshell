package log

import (
	"testing"
)

func TestProgress(t *testing.T) {
	LoadConsole(LevelDebug)
	DebugF("progress:%v", " debug")
}
