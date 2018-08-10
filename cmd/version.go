package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var version = "v2.2.2"

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

func UserAgent() string {
	return fmt.Sprintf("QShell/%s (%s; %s; %s)", version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
