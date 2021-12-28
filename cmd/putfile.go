package cmd

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload/operations"
	"github.com/spf13/cobra"
)

var formUploadCmdBuiler = func() *cobra.Command {
	info := operations.FormUploadInfo{}
	cmd := &cobra.Command{
		Use:   "fput <Bucket> <Key> <LocalFile>",
		Short: "Form upload a local file",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.FilePath = args[2]
			}
			operations.FormUpload(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "storage type")
	//cmd.Flags().IntVarP(&info.w, "worker", "c", 16, "worker count")
	cmd.Flags().StringVarP(&info.Host, "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&info.CallbackUrls, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.CallbackHost, "callback-host", "T", "", "upload callback host")

	return cmd
}

var resumeUploadCmdBuilder = func() *cobra.Command {
	info := operations.ResumeUploadInfo{}
	cmd := &cobra.Command{
		Use:   "rput <Bucket> <Key> <LocalFile>",
		Short: "Resumable upload a local file",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				info.Bucket = args[0]
				info.Key = args[1]
				info.FilePath = args[2]
			}
			operations.ResumeUpload(info)
		},
	}
	cmd.Flags().BoolVarP(&info.IsResumeV2, "v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().Int64VarP(&info.ResumeV2PartSize, "v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, default 4M")
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "storage type")
	cmd.Flags().IntVarP(&info.WorkerCount, "worker", "c", 16, "worker count")
	cmd.Flags().StringVarP(&info.Host, "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&info.CallbackUrls, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.CallbackHost, "callback-host", "T", "", "upload callback host")
	return cmd
}

func init() {
	RootCmd.AddCommand(
		forbiddenCmdBuilder(),
		resumeUploadCmdBuilder(),
		)
}
