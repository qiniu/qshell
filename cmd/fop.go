package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

var preFopStatusCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.PreFopStatusInfo{}
	var cmd = &cobra.Command{
		Use:   "prefop <PersistentId>",
		Short: "Query the pfop status",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.PreFopType
			if len(args) > 0 {
				info.Id = args[0]
			}
			operations.PreFopStatus(cfg, info)
		},
	}
	return cmd
}

var preFopCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.PreFopInfo{}
	var cmd = &cobra.Command{
		Use:   "pfop <Bucket> <Key> <fopCommand>",
		Short: "issue a request to process file in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.PFopType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.Fops = args[2]
			}
			operations.PreFop(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.Pipeline, "pipeline", "p", "", "task pipeline")
	cmd.Flags().StringVarP(&info.NotifyURL, "notify-url", "u", "", "notfiy url")
	cmd.Flags().BoolVarP(&info.NotifyForce, "force", "f", false, "force execute")
	return cmd
}

func init() {
	registerLoader(fopCmdLoader)
}

func fopCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config)  {
	superCmd.AddCommand(
		preFopCmdBuilder(cfg),
		preFopStatusCmdBuilder(cfg),
	)
}