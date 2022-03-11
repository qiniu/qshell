package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/cdn/operations"
	"github.com/spf13/cobra"
)

var cdnPrefetchCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.PrefetchInfo{}
	var cmd = &cobra.Command{
		Use:   "cdnprefetch [-i <ItemListFile>]",
		Short: "Batch prefetch the urls in the url list file",
		Long:  "Batch prefetch the urls in the url list file or from stdin if ItemListFile not specified",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.CdnPrefetchType
			operations.Prefetch(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.UrlListFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVar(&info.QpsLimit, "qps", 0, "qps limit for http call, default no limit")
	cmd.Flags().IntVarP(&info.SizeLimit, "size", "s", 50, "max item-size pre commit, max is 50, default 50")

	return cmd
}

var cdnRefreshCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.RefreshInfo{}
	var cmd = &cobra.Command{
		Use:   "cdnrefresh [-i <ItemListFile>]",
		Short: "Batch refresh the cdn cache by the url list file",
		Long:  "Batch refresh the cdn cache by the url list file or from stdin if ItemListFile not specified",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.CdnRefreshType
			operations.Refresh(cfg, info)
		},
	}

	cmd.Flags().BoolVarP(&info.IsDir, "dirs", "r", false, "refresh directory")
	cmd.Flags().StringVarP(&info.ItemListFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVar(&info.QpsLimit, "qps", 0, "qps limit for http call, default no limit")
	cmd.Flags().IntVarP(&info.SizeLimit, "size", "s", 50, "max item-size pre commit, max is 50, default 50")

	return cmd
}

func init() {
	registerLoader(cdnCmdLoader)
}

func cdnCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		cdnPrefetchCmdBuilder(cfg),
		cdnRefreshCmdBuilder(cfg),
	)
}
