package main

import (
	"fmt"
	"os"

	"github.com/qiniu/qshell/v2/cmd"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = []string{"qshell", "batchstat", "testna0", "-d"}
		os.Args = []string{"qshell", "fetch", "https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.2.html", "testna0", "-k", "06.html", "-d"}
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
