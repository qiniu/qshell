package main

import (
	"fmt"
	"os"

	"github.com/qiniu/qshell/v2/cmd"
)

func main() {
	fmt.Printf("cmd args:%v", os.Args)

	if len(os.Args) < 2 {
		os.Args = []string{"qshell", "batchstat", "testna0", "-d"}
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
