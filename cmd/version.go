package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// the version will be injected when publishing
var version = "UNSTABLE"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, params []string) {
		fmt.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

// 生成客户端代理名称
func UserAgent() string {
	return fmt.Sprintf("QShell/%s (%s; %s; %s)", version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
