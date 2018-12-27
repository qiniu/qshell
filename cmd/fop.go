package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"os"
)

var (
	prefopCmd = &cobra.Command{
		Use:   "prefop <PersistentId>",
		Short: "Query the pfop status",
		Args:  cobra.ExactArgs(1),
		Run:   Prefop,
	}

	fopCmd = &cobra.Command{
		Use:   "pfop <Bucket> <Key> <fopCommand>",
		Short: "issue a request to process file in bucket",
		Args:  cobra.ExactArgs(3),
		Run:   Fop,
	}
)

var pipeline string

func init() {
	fopCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "task pipeline")
	RootCmd.AddCommand(prefopCmd, fopCmd)
}

// 【prefop】根据persistentId查询异步处理的进度, 处理结果
func Prefop(cmd *cobra.Command, params []string) {
	persistentId := params[0]

	fopRet, err := iqshell.Prefop(persistentId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prefop error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	} else {
		fmt.Println(fopRet.String())
	}
}

// 【pfop】 提交异步处理请求
func Fop(cmd *cobra.Command, params []string) {
	bucket, key, fops := params[0], params[1], params[2]

	persistengId, err := iqshell.Pfop(bucket, key, fops, pipeline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prefop error: %v\n", err)
		os.Exit(iqshell.STATUS_ERROR)
	}
	fmt.Println(persistengId)
}
