package cmd

import "os"

func setOsArgsAndRun(args []string) {
	os.Args = args
	Execute()
}
