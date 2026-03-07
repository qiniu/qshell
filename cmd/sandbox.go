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
		Short:   "Manage sandboxes (alias: sbx)",
		Example: `  # View sandbox subcommands
  qshell sandbox -h
  qshell sbx -h

  # Create a sandbox from a template
  qshell sandbox create my-template
  qshell sbx cr my-template

  # List running sandboxes
  qshell sandbox list
  qshell sbx ls

  # Connect to a sandbox
  qshell sandbox connect sb-xxxxxxxxxxxx
  qshell sbx cn sb-xxxxxxxxxxxx`,
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
		Short:   "List sandboxes (alias: ls)",
		Example: `  # List running sandboxes
  qshell sandbox list
  qshell sbx ls

  # Filter by state and limit
  qshell sandbox list --state running,paused --limit 10
  qshell sbx ls -s running,paused -l 10

  # Filter by metadata
  qshell sandbox list -m env=prod,team=backend
  qshell sbx ls -m env=prod,team=backend

  # Output as JSON
  qshell sandbox list -f json
  qshell sbx ls -f json`,
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
		Short:   "Create a sandbox and connect to its terminal (alias: cr)",
		Example: `  # Create a sandbox from a template
  qshell sandbox create my-template
  qshell sbx cr my-template

  # Create with a timeout (seconds)
  qshell sandbox create my-template --timeout 300
  qshell sbx cr my-template -t 300

  # Create with metadata
  qshell sandbox create my-template -m env=dev,team=backend
  qshell sbx cr my-template -m env=dev,team=backend`,
		Args: cobra.MaximumNArgs(1),
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
	cmd.Flags().Int32VarP(&info.Timeout, "timeout", "t", 0, "sandbox timeout in seconds")
	cmd.Flags().StringVarP(&info.Metadata, "metadata", "m", "", "metadata key=value pairs (comma-separated)")
	return cmd
}

var sandboxConnectCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect <sandboxID>",
		Aliases: []string{"cn"},
		Short:   "Connect to an existing sandbox terminal (alias: cn)",
		Example: `  # Connect to a sandbox by ID
  qshell sandbox connect sb-xxxxxxxxxxxx
  qshell sbx cn sb-xxxxxxxxxxxx`,
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
		Short:   "Kill one or more sandboxes (alias: kl)",
		Example: `  # Kill a single sandbox
  qshell sandbox kill sb-xxxxxxxxxxxx
  qshell sbx kl sb-xxxxxxxxxxxx

  # Kill multiple sandboxes
  qshell sandbox kill sb-aaa sb-bbb sb-ccc
  qshell sbx kl sb-aaa sb-bbb sb-ccc

  # Kill all running sandboxes
  qshell sandbox kill --all
  qshell sbx kl -a

  # Kill all paused sandboxes with specific metadata
  qshell sandbox kill --all -s paused -m env=dev
  qshell sbx kl -a -s paused -m env=dev`,
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
		Short:   "View sandbox logs (alias: lg)",
		Example: `  # View logs
  qshell sandbox logs sb-xxxxxxxxxxxx
  qshell sbx lg sb-xxxxxxxxxxxx

  # Filter by level (WARN and above)
  qshell sandbox logs sb-xxxxxxxxxxxx --level WARN
  qshell sbx lg sb-xxxxxxxxxxxx --level WARN

  # Stream logs in follow mode
  qshell sandbox logs sb-xxxxxxxxxxxx -f
  qshell sbx lg sb-xxxxxxxxxxxx -f

  # Filter by logger prefix
  qshell sandbox logs sb-xxxxxxxxxxxx --loggers envd,process
  qshell sbx lg sb-xxxxxxxxxxxx --loggers envd,process

  # Output as JSON
  qshell sandbox logs sb-xxxxxxxxxxxx --format json
  qshell sbx lg sb-xxxxxxxxxxxx --format json`,
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
		Short:   "View sandbox resource metrics (alias: mt)",
		Example: `  # View current metrics
  qshell sandbox metrics sb-xxxxxxxxxxxx
  qshell sbx mt sb-xxxxxxxxxxxx

  # Stream metrics in follow mode
  qshell sandbox metrics sb-xxxxxxxxxxxx -f
  qshell sbx mt sb-xxxxxxxxxxxx -f

  # Output as JSON
  qshell sandbox metrics sb-xxxxxxxxxxxx --format json
  qshell sbx mt sb-xxxxxxxxxxxx --format json`,
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
