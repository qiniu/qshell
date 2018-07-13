package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qiniu/rpc"
	"github.com/tonycai653/iqshell/qshell"
	"os"
)

var prefopCmd = &cobra.Command{
	Use:   "prefop <PersistentId>",
	Short: "Query the pfop status",
	Args:  cobra.ExactArgs(1),
	Run:   Prefop,
}

func init() {
	RootCmd.AddCommand(prefopCmd)
}

func Prefop(cmd *cobra.Command, params []string) {
	persistentId := params[0]
	fopRet := qshell.FopRet{}
	err := qshell.Prefop(persistentId, &fopRet)
	if err != nil {
		if v, ok := err.(*rpc.ErrorInfo); ok {
			fmt.Println("Prefop error,", v.Code, v.Err)
		} else {
			fmt.Println("Prefop error,", err)
		}
		os.Exit(qshell.STATUS_ERROR)
	} else {
		fmt.Println(fopRet.String())
	}
}
