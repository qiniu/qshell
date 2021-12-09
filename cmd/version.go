package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

func versionCmdBuilder() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, params []string) {
			fmt.Println(data.Version)
		},
	}
	return cmd
}

func testCmdBuilder() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "test",
		Short: "just test",
		Run: func(cmd *cobra.Command, params []string) {
			test()
		},
	}
	return cmd
}

func init() {
	RootCmd.AddCommand(versionCmdBuilder(), testCmdBuilder())
}

type Status struct {
	isCancel bool
}

func test() {
	status := &Status{
		isCancel: false,
	}

	c := make(chan int, 1)
	go func(c chan int) {
		for i := 0; i < 20; i++ {
			c <- i
			time.Sleep(time.Second)
		}
	}(c)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	go func(c chan int) {
		si := <-s
		log.ErrorF("Got signal:", si)
		status.isCancel = true
	}(c)

	for i := range c {
		log.AlertF("=== out:%d", i)
		if status.isCancel {
			break
		}
	}

	for i := 0; i < 3; i++ {
		log.AlertF("=== completed ===")
		time.Sleep(time.Second)
	}
}
