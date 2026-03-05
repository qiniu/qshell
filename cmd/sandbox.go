package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/sandbox/operations"
)

var sandboxCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sandbox",
		Aliases: []string{"sbx"},
		Short:   "Manage sandboxes",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxType
			docs.ShowCmdDocument(docs.SandboxType)
		},
	}
	return cmd
}

var sandboxListCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ListInfo{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List sandboxes",
		Example: `qshell sandbox list --state running --limit 10`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxListType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.List(info)
		},
	}
	cmd.Flags().StringVarP(&info.State, "state", "s", "", "filter by state (comma-separated: running,paused). Defaults to running")
	cmd.Flags().StringVarP(&info.Metadata, "metadata", "m", "", "filter by metadata (key1=value1,key2=value2)")
	cmd.Flags().Int32VarP(&info.Limit, "limit", "l", 0, "maximum number of sandboxes to return")
	cmd.Flags().StringVarP(&info.Format, "format", "f", "pretty", "output format: pretty or json")
	return cmd
}

var sandboxCreateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.CreateInfo{}
	cmd := &cobra.Command{
		Use:     "create [template]",
		Aliases: []string{"cr"},
		Short:   "Create a sandbox and connect to its terminal",
		Example: `qshell sandbox create my-template`,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxCreateType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) > 0 {
				info.TemplateID = args[0]
			}
			operations.Create(info)
		},
	}
	return cmd
}

var sandboxConnectCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect <sandboxID>",
		Aliases: []string{"cn"},
		Short:   "Connect to an existing sandbox terminal",
		Example: `qshell sandbox connect sb-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxConnectType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 1 {
				_ = cmd.Usage()
				return
			}
			operations.Connect(operations.ConnectInfo{
				SandboxID: args[0],
			})
		},
	}
	return cmd
}

var sandboxKillCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.KillInfo{}
	cmd := &cobra.Command{
		Use:     "kill [sandboxIDs...]",
		Aliases: []string{"kl"},
		Short:   "Kill one or more sandboxes",
		Example: `qshell sandbox kill sb-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxKillType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			info.SandboxIDs = args
			operations.Kill(info)
		},
	}
	cmd.Flags().BoolVarP(&info.All, "all", "a", false, "kill all sandboxes")
	cmd.Flags().StringVarP(&info.State, "state", "s", "", "filter by state when using --all (comma-separated: running,paused). Defaults to running")
	cmd.Flags().StringVarP(&info.Metadata, "metadata", "m", "", "filter by metadata when using --all (key1=value1,key2=value2)")
	return cmd
}

var sandboxLogsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.LogsInfo{}
	cmd := &cobra.Command{
		Use:     "logs <sandboxID>",
		Aliases: []string{"lg"},
		Short:   "View sandbox logs",
		Example: `qshell sandbox logs sb-xxxxxxxxxxxx --level INFO`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxLogsType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 1 {
				_ = cmd.Usage()
				return
			}
			info.SandboxID = args[0]
			operations.Logs(info)
		},
	}
	cmd.Flags().StringVar(&info.Level, "level", "INFO", "filter by log level (DEBUG, INFO, WARN, ERROR). Higher levels are also shown")
	cmd.Flags().Int32Var(&info.Limit, "limit", 0, "maximum number of log entries to return")
	cmd.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	cmd.Flags().BoolVarP(&info.Follow, "follow", "f", false, "keep streaming logs until the sandbox is closed")
	cmd.Flags().StringVar(&info.Loggers, "loggers", "", "filter logs by loggers (comma-separated prefixes)")
	return cmd
}

var sandboxMetricsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.MetricsInfo{}
	cmd := &cobra.Command{
		Use:     "metrics <sandboxID>",
		Aliases: []string{"mt"},
		Short:   "View sandbox resource metrics",
		Example: `qshell sandbox metrics sb-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxMetricsType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 1 {
				_ = cmd.Usage()
				return
			}
			info.SandboxID = args[0]
			operations.Metrics(info)
		},
	}
	cmd.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	cmd.Flags().BoolVarP(&info.Follow, "follow", "f", false, "keep streaming metrics until the sandbox is closed")
	return cmd
}

func init() {
	registerLoader(sandboxCmdLoader)
}

func sandboxCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	sandboxCmd := sandboxCmdBuilder(cfg)
	sandboxCmd.AddCommand(
		sandboxListCmdBuilder(cfg),
		sandboxCreateCmdBuilder(cfg),
		sandboxConnectCmdBuilder(cfg),
		sandboxKillCmdBuilder(cfg),
		sandboxLogsCmdBuilder(cfg),
		sandboxMetricsCmdBuilder(cfg),
	)

	// Add template as a subcommand of sandbox
	templateCmdLoader(sandboxCmd, cfg)

	superCmd.AddCommand(sandboxCmd)
}
