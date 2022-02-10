package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8/operations"
	"github.com/spf13/cobra"
)

var m3u8ReplaceDomainCmdBuilder = func() *cobra.Command {
	var info = operations.ReplaceDomainInfo{}
	var cmd = &cobra.Command{
		Use:   "m3u8replace <Bucket> <M3u8Key> [<NewDomain>]",
		Short: "Replace m3u8 domain in the playlist",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			if len(args) == 3 {
				info.NewDomain = args[2]
			}
			prepare(cmd, nil)
			operations.ReplaceDomain(info)
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
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				info.Bucket = args[0]
				info.Key = args[1]
			}
			prepare(cmd, nil)
			operations.Delete(info)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(m3u8ReplaceDomainCmdBuilder(), m3u8DeleteCmdBuilder())
}
