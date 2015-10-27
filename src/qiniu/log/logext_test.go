package log

import (
	"testing"
)

func TestLog(t *testing.T) {

	SetOutputLevel(Ldebug)

	Debugf("Debug: foo\n")
	Debug("Debug: foo")

	Infof("Info: foo\n")
	Info("Info: foo")

	Warnf("Warn: foo\n")
	Warn("Warn: foo")

	Errorf("Error: foo\n")
	Error("Error: foo")

	SetOutputLevel(Linfo)

	Debugf("Debug: foo\n")
	Debug("Debug: foo")

	Infof("Info: foo\n")
	Info("Info: foo")

	Warnf("Warn: foo\n")
	Warn("Warn: foo")

	Errorf("Error: foo\n")
	Error("Error: foo")
}
