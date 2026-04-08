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

  # Create in detached mode (no terminal, sandbox stays alive)
  qshell sandbox create my-template -t 300 --detach
  qshell sbx cr my-template -t 300 --detach

  # Create with environment variables
  qshell sandbox create my-template -e FOO=bar -e BAZ=qux
  qshell sbx cr my-template -e FOO=bar -e BAZ=qux

  # Create with auto-pause (pause instead of kill on timeout)
  qshell sandbox create my-template -t 300 --auto-pause
  qshell sbx cr my-template -t 300 --auto-pause

  # Create with metadata
  qshell sandbox create my-template -m env=dev,team=backend
  qshell sbx cr my-template -m env=dev,team=backend

  # Create with injection rules
  qshell sandbox create my-template --injection-rule rule-openai --injection-rule rule-http
  qshell sbx cr my-template --injection-rule rule-openai --injection-rule rule-http

  # Create with inline injections
  qshell sandbox create my-template --inline-injection 'type=openai,api-key=sk-xxx' --inline-injection 'type=http,base-url=https://api.example.com,headers=Authorization=Bearer token;X-Env=prod'
  qshell sbx cr my-template --inline-injection 'type=openai,api-key=sk-xxx'`,
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
	cmd.Flags().BoolVar(&info.Detach, "detach", false, "create sandbox without connecting terminal (sandbox stays alive until timeout)")
	cmd.Flags().StringVarP(&info.Metadata, "metadata", "m", "", "metadata key=value pairs (comma-separated)")
	cmd.Flags().StringArrayVarP(&info.EnvVars, "env-var", "e", nil, "environment variables (KEY=VALUE, can be specified multiple times)")
	cmd.Flags().BoolVar(&info.AutoPause, "auto-pause", false, "automatically pause sandbox when timeout expires (instead of killing)")
	cmd.Flags().StringArrayVar(&info.InjectionRuleID, "injection-rule", nil, "injection rule IDs to apply when creating the sandbox (can be specified multiple times)")
	cmd.Flags().StringArrayVar(&info.InlineInjection, "inline-injection", nil, "inline injection spec to apply when creating the sandbox (can be specified multiple times, format: type=<type>,api-key=<key>,base-url=<url>,headers=<k1=v1;k2=v2>)")
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

var sandboxPauseCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.PauseInfo{}
	cmd := &cobra.Command{
		Use:     "pause [sandboxIDs...]",
		Aliases: []string{"ps"},
		Short:   "Pause one or more sandboxes (alias: ps)",
		Example: `  # Pause a single sandbox
  qshell sandbox pause sb-xxxxxxxxxxxx
  qshell sbx ps sb-xxxxxxxxxxxx

  # Pause multiple sandboxes
  qshell sandbox pause sb-aaa sb-bbb sb-ccc
  qshell sbx ps sb-aaa sb-bbb sb-ccc

  # Pause all running sandboxes
  qshell sandbox pause --all
  qshell sbx ps -a

  # Pause all with specific metadata
  qshell sandbox pause --all -m env=dev
  qshell sbx ps -a -m env=dev`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxPauseType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			info.SandboxIDs = args
			operations.Pause(info)
		},
	}
	cmd.Flags().BoolVarP(&info.All, "all", "a", false, "pause all sandboxes")
	cmd.Flags().StringVarP(&info.State, "state", "s", "", "filter by state when using --all (comma-separated: running,paused). Defaults to running")
	cmd.Flags().StringVarP(&info.Metadata, "metadata", "m", "", "filter by metadata when using --all (key1=value1,key2=value2)")
	return cmd
}

var sandboxResumeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ResumeInfo{}
	cmd := &cobra.Command{
		Use:     "resume [sandboxIDs...]",
		Aliases: []string{"rs"},
		Short:   "Resume one or more paused sandboxes (alias: rs)",
		Example: `  # Resume a paused sandbox
  qshell sandbox resume sb-xxxxxxxxxxxx
  qshell sbx rs sb-xxxxxxxxxxxx

  # Resume multiple sandboxes
  qshell sandbox resume sb-aaa sb-bbb sb-ccc
  qshell sbx rs sb-aaa sb-bbb sb-ccc

  # Resume all paused sandboxes
  qshell sandbox resume --all
  qshell sbx rs -a

  # Resume all with specific metadata
  qshell sandbox resume --all -m env=staging
  qshell sbx rs -a -m env=staging`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxResumeType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			info.SandboxIDs = args
			operations.Resume(info)
		},
	}
	cmd.Flags().BoolVarP(&info.All, "all", "a", false, "resume all paused sandboxes")
	cmd.Flags().StringVarP(&info.Metadata, "metadata", "m", "", "filter by metadata when using --all (key1=value1,key2=value2)")
	return cmd
}

var sandboxExecCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ExecInfo{}
	cmd := &cobra.Command{
		Use:     "exec <sandboxID> -- <command...>",
		Aliases: []string{"ex"},
		Short:   "Execute a command in a sandbox (alias: ex)",
		Example: `  # Run a command in a sandbox
  qshell sandbox exec sb-xxxxxxxxxxxx -- ls -la
  qshell sbx ex sb-xxxxxxxxxxxx -- ls -la

  # Pipe stdin to a command
  echo "hello world" | qshell sbx ex sb-xxxxxxxxxxxx -- cat
  cat file.txt | qshell sbx ex sb-xxxxxxxxxxxx -- wc -l

  # Run in background (print PID and return)
  qshell sandbox exec sb-xxxxxxxxxxxx -b -- python server.py
  qshell sbx ex sb-xxxxxxxxxxxx -b -- python server.py

  # Specify working directory and user
  qshell sandbox exec sb-xxxxxxxxxxxx -c /app -u root -- npm install

  # Set environment variables
  qshell sandbox exec sb-xxxxxxxxxxxx -e PORT=3000 -e NODE_ENV=production -- node app.js`,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: false,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxExecType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) < 1 {
				_ = cmd.Usage()
				return
			}
			info.SandboxID = args[0]
			// args after the sandbox ID are the command (cobra handles -- separator)
			if dash := cmd.ArgsLenAtDash(); dash >= 0 {
				info.Command = args[dash:]
			} else if len(args) > 1 {
				info.Command = args[1:]
			}
			operations.Exec(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Background, "background", "b", false, "run command in background (print PID and return)")
	cmd.Flags().StringVarP(&info.Cwd, "cwd", "c", "", "working directory for the command")
	cmd.Flags().StringVarP(&info.User, "user", "u", "", "user to run the command as")
	cmd.Flags().StringArrayVarP(&info.Envs, "env", "e", nil, "environment variables (KEY=VALUE, can be specified multiple times)")
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
		sandboxPauseCmdBuilder(cfg),
		sandboxResumeCmdBuilder(cfg),
		sandboxExecCmdBuilder(cfg),
		sandboxLogsCmdBuilder(cfg),
		sandboxMetricsCmdBuilder(cfg),
	)

	// Add template as a subcommand of sandbox
	templateCmdLoader(sandboxCmd, cfg)

	// Add injection-rule as a subcommand of sandbox
	injectionRuleCmdLoader(sandboxCmd, cfg)

	superCmd.AddCommand(sandboxCmd)
}
