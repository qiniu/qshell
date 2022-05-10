package main

import (
	"fmt"
	"os"

	"github.com/qiniu/qshell/v2/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
