package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8/operations"
	"github.com/spf13/cobra"
)

var m3u8ReplaceDomainCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ReplaceDomainInfo{}
	var cmd = &cobra.Command{
		Use:   "m3u8replace <Bucket> <M3u8Key> [<NewDomain>]",
		Short: "Replace m3u8 domain in the playlist",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.M3u8ReplaceType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.NewDomain = args[2]
			}
			operations.ReplaceDomain(cfg, info)
		},
	}

	cmd.Flags().BoolVarP(&info.RemoveSparePreSlash, "remove-spare-pre-slash", "r", true, "remove spare prefix slash(/) , only keep one slash if ts path has prefix / ")

	return cmd
}

var m3u8DeleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.DeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "m3u8delete <Bucket> <M3u8Key>",
		Short: "Delete m3u8 playlist and the slices it references",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.M3u8DeleteType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.Delete(cfg, info)
		},
	}
	return cmd
}

func init() {
	registerLoader(m3u8CmdLoader)
}

func m3u8CmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		m3u8ReplaceDomainCmdBuilder(cfg),
		m3u8DeleteCmdBuilder(cfg),
	)
}
