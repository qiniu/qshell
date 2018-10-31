package main

import (
	"github.com/tonycai653/iqshell/cmd"
	/*
		"bufio"
		"fmt"
	*/
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	/*
		fh, err := os.Open("/Users/caijiaqiang/projects/go/iqshell/keyList.txt")
		if err != nil {
			os.Exit(1)
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	*/

}
