package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/operations"
)

var templateTagsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.TagsInfo{}
	command := &cobra.Command{Use: "tags <templateID>", Short: "List template tags and assigned builds", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateTagsType
			if !iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{}) {
				return nil
			}
			info.TemplateID = args[0]
			return operations.Tags(cmd.Context(), info)
		}}
	command.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	return command
}
