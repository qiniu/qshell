package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/sandbox/operations"
)

var resourceCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resource",
		Aliases: []string{"resources"},
		Short:   "Manage resources mounted in sandboxes (alias: resources)",
		Args:    cobra.NoArgs,
		Example: `  # View resource subcommands
  qshell sandbox resource -h
  qshell sbx resource -h

  # List resources mounted in a sandbox
  qshell sandbox resource list sb-xxxxxxxxxxxx
  qshell sbx resource ls sb-xxxxxxxxxxxx

  # Update a Git repository resource token
  qshell sandbox resource update sb-xxxxxxxxxxxx res-xxxxxxxxxxxx --token ghp-xxx
  qshell sbx resource up sb-xxxxxxxxxxxx res-xxxxxxxxxxxx -t ghp-xxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxResourceType
			docs.ShowCmdDocument(docs.SandboxResourceType)
		},
	}
	return cmd
}

var resourceListCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ResourceListInfo{Format: "pretty"}
	cmd := &cobra.Command{
		Use:     "list <sandboxID>",
		Aliases: []string{"ls"},
		Short:   "List resources mounted in a sandbox (alias: ls)",
		Args:    cobra.ExactArgs(1),
		Example: `  # List resources mounted in a sandbox
  qshell sandbox resource list sb-xxxxxxxxxxxx
  qshell sbx resource ls sb-xxxxxxxxxxxx

  # Output resource details as JSON
  qshell sandbox resource list sb-xxxxxxxxxxxx --format json`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxResourceListType
			if !iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{}) {
				return
			}
			info.SandboxID = args[0]
			operations.ListResources(info)
		},
	}
	cmd.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	return cmd
}

var resourceUpdateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ResourceUpdateInfo{}
	cmd := &cobra.Command{
		Use:     "update <sandboxID> <resourceID>",
		Aliases: []string{"up", "update-token"},
		Short:   "Update a Git repository resource token (alias: up)",
		Args:    cobra.ExactArgs(2),
		Example: `  # Update a GitHub repository resource token
  qshell sandbox resource update sb-xxxxxxxxxxxx res-xxxxxxxxxxxx --token ghp-xxx
  qshell sbx resource up sb-xxxxxxxxxxxx res-xxxxxxxxxxxx -t ghp-xxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxResourceUpdateType
			if !iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{}) {
				return
			}
			info.SandboxID = args[0]
			info.ResourceID = args[1]
			operations.UpdateResourceToken(info)
		},
	}
	cmd.Flags().StringVarP(&info.Token, "token", "t", "", "new Git repository authorization token")
	cmd.Flags().StringVar(&info.Token, "authorization-token", "", "alias for --token")
	return cmd
}

func resourceCmdLoader(parentCmd *cobra.Command, cfg *iqshell.Config) {
	resourceCmd := resourceCmdBuilder(cfg)
	resourceCmd.AddCommand(
		resourceListCmdBuilder(cfg),
		resourceUpdateCmdBuilder(cfg),
	)
	parentCmd.AddCommand(resourceCmd)
}
