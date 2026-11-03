package cmd

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/sandbox/operations"
)

var sandboxSnapshotCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	command := &cobra.Command{Use: "snapshot", Short: "Manage persistent sandbox snapshots", Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxSnapshotType
			docs.ShowCmdDocument(docs.SandboxSnapshotType)
		},
	}
	createInfo := operations.SnapshotCreateInfo{}
	create := &cobra.Command{Use: "create <sandboxID>", Short: "Create a persistent snapshot", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.CmdCfg.CmdId = docs.SandboxSnapshotCreateType
			if !iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{}) {
				return nil
			}
			createInfo.SandboxID = args[0]
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer cancel()
			return operations.SnapshotCreate(ctx, createInfo)
		},
	}
	create.Flags().StringVar(&createInfo.Name, "name", "", "snapshot alias, optionally including a tag")
	create.Flags().BoolVar(&createInfo.Memory, "memory", true, "include memory; false creates a filesystem-only snapshot that cold-boots")
	create.Flags().StringVar(&createInfo.IdempotencyKey, "idempotency-key", "", "required stable key; reuse it with the same sandbox and request body after an uncertain outcome")
	create.Flags().StringVar(&createInfo.Format, "format", "pretty", "output format: pretty or json")
	_ = create.MarkFlagRequired("idempotency-key")
	listInfo := operations.SnapshotListInfo{}
	list := &cobra.Command{Use: "list", Short: "List a page of persistent snapshots", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.CmdCfg.CmdId = docs.SandboxSnapshotListType
			if !iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{}) {
				return nil
			}
			return operations.SnapshotList(cmd.Context(), listInfo)
		},
	}
	list.Flags().StringVar(&listInfo.Name, "name", "", "filter by snapshot name or template ID")
	list.Flags().StringVar(&listInfo.SandboxID, "sandbox-id", "", "filter by source sandbox ID")
	list.Flags().Int32Var(&listInfo.Limit, "limit", 0, "maximum items per page; omitted uses the server default")
	list.Flags().StringVar(&listInfo.NextToken, "next-token", "", "pagination cursor returned by the previous page")
	list.Flags().StringVar(&listInfo.Format, "format", "pretty", "output format: pretty or json")
	command.AddCommand(create, list)
	return command
}
