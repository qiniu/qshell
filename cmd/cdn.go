package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/cdn/operations"
	"github.com/spf13/cobra"
)

var cdnPrefetchCmdBuilder = func() *cobra.Command {
	var info = operations.PrefetchInfo{}
	var cmd = &cobra.Command{
		Use:   "cdnprefetch [-i <ItemListFile>]",
		Short: "Batch prefetch the urls in the url list file",
		Long:  "Batch prefetch the urls in the url list file or from stdin if ItemListFile not specified",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			loadConfig()
			operations.Prefetch(info)
		},
	}

	cmd.Flags().StringVarP(&info.UrlListFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVar(&info.QpsLimit, "qps", 0, "qps limit for http call")
	cmd.Flags().IntVarP(&info.SizeLimit, "size", "s", 0, "max item-size pre commit")

	return cmd
}

var cdnRefreshCmdBuilder = func() *cobra.Command {
	var info = operations.RefreshInfo{}
	var cmd = &cobra.Command{
		Use:   "cdnrefresh [-i <ItemListFile>]",
		Short: "Batch refresh the cdn cache by the url list file",
		Long:  "Batch refresh the cdn cache by the url list file or from stdin if ItemListFile not specified",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			loadConfig()
			operations.Refresh(info)
		},
	}

	cmd.Flags().BoolVarP(&info.IsDir, "dirs", "r", false, "refresh directory")
	cmd.Flags().StringVarP(&info.ItemListFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVar(&info.QpsLimit, "qps", 0, "qps limit for http call")
	cmd.Flags().IntVarP(&info.SizeLimit, "size", "s", 0, "max item-size pre commit")

	return cmd
}

func init() {
	RootCmd.AddCommand(
		cdnPrefetchCmdBuilder(),
		cdnRefreshCmdBuilder(),
	)
}
