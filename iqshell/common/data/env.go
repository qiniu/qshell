package data

import "os"

const (
	EnvKeyTest = "qshell_test"
)

func IsTestMode() bool {
	return os.Getenv(EnvKeyTest) == TrueString
}

func SetTestMode() *CodeError {
	return ConvertError(os.Setenv(EnvKeyTest, TrueString))
}

