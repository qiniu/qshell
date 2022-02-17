package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8/operations"
	"github.com/spf13/cobra"
)

var m3u8ReplaceDomainCmdBuilder = func() *cobra.Command {
	var info = operations.ReplaceDomainInfo{}
	var cmd = &cobra.Command{
		Use:   "m3u8replace <Bucket> <M3u8Key> [<NewDomain>]",
		Short: "Replace m3u8 domain in the playlist",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.M3u8ReplaceType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.NewDomain = args[2]
			}
			if prepare(cmd, &info) {
				operations.ReplaceDomain(info)
			}
		},
	}

	cmd.Flags().BoolVarP(&info.RemoveSparePreSlash, "remove-spare-pre-slash", "r", true, "remove spare prefix slash(/) , only keep one slash if ts path has prefix / ")

	return cmd
}

var m3u8DeleteCmdBuilder = func() *cobra.Command {
	var info = operations.DeleteInfo{}
	var cmd = &cobra.Command{
		Use:   "m3u8delete <Bucket> <M3u8Key>",
		Short: "Delete m3u8 playlist and the slices it references",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.M3u8DeleteType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if prepare(cmd, &info) {
				operations.Delete(info)
			}
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(m3u8ReplaceDomainCmdBuilder(), m3u8DeleteCmdBuilder())
}
