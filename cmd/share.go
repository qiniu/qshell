package cmd

import (
	"strings"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/operations"
	"github.com/spf13/cobra"
)

var createShareCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.CreateShareInfo{}
	var cmd = &cobra.Command{
		Use:   "create-share [kodo://]<Bucket>/<Prefix>",
		Short: "Share the specified directory",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.CreateShareType
			if len(args) > 0 {
				url := args[0]
				url = strings.TrimPrefix(url, "kodo://")
				parts := strings.SplitN(url, "/", 2)
				if len(parts) > 0 {
					info.Bucket = parts[0]
				}
				if len(parts) > 1 {
					info.Prefix = parts[1]
				}
				info.Permission = "READONLY"
			}
			operations.CreateShare(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.ExtractCode, "extract-code", "", "", "Specify extract code for the share, must consist of 6 alphabets and numbers")
	cmd.Flags().Int64VarP(&info.DurationSeconds, "validity-period", "", 900, "Specify validity period in seconds, must be longer than 15 mintutes and shorter than 2 hours")
	cmd.Flags().StringVarP(&info.OutputPath, "output", "", "", "Specify path to save output")
	return cmd
}

var listShareCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ListShareInfo{}
	var cmd = &cobra.Command{
		Use:   "share-ls <Link>",
		Short: "List shared files",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ShareLsType
			if len(args) > 0 {
				info.LinkURL = args[0]
			}
			operations.ListShare(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.ExtractCode, "extract-code", "", "", "Specify extract code for the share, must consist of 6 alphabets and numbers")
	cmd.Flags().StringVarP(&info.Prefix, "prefix", "p", "", "list the shared directory with this prefix if set")
	cmd.Flags().Int64VarP(&info.Limit, "limit", "n", 0, "max count of items to output")
	cmd.Flags().StringVarP(&info.Marker, "marker", "m", "", "list marker")
	return cmd
}

var copyShareCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.CopyShareInfo{}
	var cmd = &cobra.Command{
		Use:   "share-cp --link=<Link> --to=<ToPath>",
		Short: "Copy shared files",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ShareCpType
			if len(args) > 0 {
				info.LinkURL = args[0]
			}
			operations.CopyShare(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.ExtractCode, "extract-code", "", "", "Specify extract code for the share, must consist of 6 alphabets and numbers")
	cmd.Flags().StringVarP(&info.FromPath, "from", "", "", "copy from path")
	cmd.Flags().StringVarP(&info.ToPath, "to", "", "", "copy to path")
	return cmd
}

func init() {
	registerLoader(shareCmdLoader)
}

func shareCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		createShareCmdBuilder(cfg),
		listShareCmdBuilder(cfg),
		copyShareCmdBuilder(cfg),
	)
}
