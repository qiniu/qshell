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

	cmd.Flags().StringVarP(&info.Bucket, "bucket", "", "", "the bucket that the fop file in, public cloud is optional, private cloud must be set when no api host is configured.")

	return cmd
}

var preFopCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.PreFopInfo{}
	var cmd = &cobra.Command{
		Use:   "pfop <Bucket> <Key> <fopCommand>",
		Short: "Issue a request to process file in bucket",
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

	cmd.Flags().BoolVarP(&info.NotifyForce, "force", "y", false, "force execute")
	cmd.Flags().BoolVarP(&info.NotifyForce, "force-old", "f", false, "force execute, deprecated")
	_ = cmd.Flags().MarkDeprecated("force-old", "use --force instead")

	return cmd
}

func init() {
	registerLoader(fopCmdLoader)
}

func fopCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		preFopCmdBuilder(cfg),
		preFopStatusCmdBuilder(cfg),
	)
}
