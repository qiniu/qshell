package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/version/operations"
	"github.com/spf13/cobra"
)

func versionCmdBuilder(cfg *iqshell.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, params []string) {
			cfg.CmdCfg.CmdId = docs.VersionType
			operations.Version(cfg, operations.VersionInfo{})
		},
	}
	return cmd
}

func init() {
	registerLoader(versionCmdLoader)
}

func versionCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		versionCmdBuilder(cfg),
	)
}
