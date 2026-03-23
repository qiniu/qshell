package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/transform_rule/operations"
)

var transformRuleCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transform-rule",
		Aliases: []string{"tr"},
		Short:   "Manage sandbox transform rules (alias: tr)",
		Example: `  # View transform-rule subcommands
  qshell sandbox transform-rule -h
  qshell sbx tr -h

  # List all transform rules
  qshell sandbox transform-rule list
  qshell sbx tr ls

  # Create a transform rule
  qshell sandbox transform-rule create --name my-rule --hosts api.example.com
  qshell sbx tr cr --name my-rule --hosts api.example.com

  # Get transform rule details
  qshell sandbox transform-rule get rule-xxxxxxxxxxxx
  qshell sbx tr gt rule-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTransformRuleType
			docs.ShowCmdDocument(docs.SandboxTransformRuleType)
		},
	}
	return cmd
}

var transformRuleListCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ListInfo{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List transform rules (alias: ls)",
		Example: `  # List all transform rules
  qshell sandbox transform-rule list
  qshell sbx tr ls

  # Output as JSON
  qshell sandbox transform-rule list --format json
  qshell sbx tr ls --format json`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTransformRuleListType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.List(info)
		},
	}
	cmd.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	return cmd
}

var transformRuleGetCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <ruleID>",
		Aliases: []string{"gt"},
		Short:   "Get transform rule details (alias: gt)",
		Example: `  # Get transform rule details
  qshell sandbox transform-rule get rule-xxxxxxxxxxxx
  qshell sbx tr gt rule-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTransformRuleGetType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 1 {
				_ = cmd.Usage()
				return
			}
			operations.Get(operations.GetInfo{
				RuleID: args[0],
			})
		},
	}
	return cmd
}

var transformRuleCreateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.CreateInfo{}
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"cr"},
		Short:   "Create a transform rule (alias: cr)",
		Example: `  # Create a basic transform rule
  qshell sandbox transform-rule create --name my-rule --hosts api.example.com
  qshell sbx tr cr --name my-rule --hosts api.example.com

  # Create with headers replacement
  qshell sandbox transform-rule create --name api-auth --hosts api.example.com --headers "Authorization=Bearer token123"

  # Create with headers and queries
  qshell sandbox transform-rule create --name full-rule --hosts api.example.com,cdn.example.com --headers "Authorization=Bearer xxx" --queries "token=abc"`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTransformRuleCreateType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.Create(info)
		},
	}
	cmd.Flags().StringVar(&info.Name, "name", "", "rule name (required, unique per user)")
	cmd.Flags().StringVar(&info.Hosts, "hosts", "", "match hosts (comma-separated, e.g. api.example.com,cdn.example.com)")
	cmd.Flags().StringVar(&info.Headers, "headers", "", "replacement headers (comma-separated key=value pairs)")
	cmd.Flags().StringVar(&info.Queries, "queries", "", "replacement query parameters (comma-separated key=value pairs)")
	return cmd
}

var transformRuleUpdateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UpdateInfo{}
	cmd := &cobra.Command{
		Use:     "update <ruleID>",
		Aliases: []string{"up"},
		Short:   "Update a transform rule (alias: up)",
		Example: `  # Update rule name
  qshell sandbox transform-rule update rule-xxxxxxxxxxxx --name new-name
  qshell sbx tr up rule-xxxxxxxxxxxx --name new-name

  # Update hosts
  qshell sandbox transform-rule update rule-xxxxxxxxxxxx --hosts api.new-domain.com

  # Update multiple fields
  qshell sandbox transform-rule update rule-xxxxxxxxxxxx --name updated-rule --hosts api.example.com --headers "Authorization=Bearer newtoken"`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTransformRuleUpdateType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 1 {
				_ = cmd.Usage()
				return
			}
			info.RuleID = args[0]
			operations.Update(info)
		},
	}
	cmd.Flags().StringVar(&info.Name, "name", "", "new rule name")
	cmd.Flags().StringVar(&info.Hosts, "hosts", "", "new match hosts (comma-separated)")
	cmd.Flags().StringVar(&info.Headers, "headers", "", "new replacement headers (comma-separated key=value pairs)")
	cmd.Flags().StringVar(&info.Queries, "queries", "", "new replacement query parameters (comma-separated key=value pairs)")
	return cmd
}

var transformRuleDeleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DeleteInfo{}
	cmd := &cobra.Command{
		Use:     "delete [ruleIDs...]",
		Aliases: []string{"dl"},
		Short:   "Delete one or more transform rules (alias: dl)",
		Example: `  # Delete a single rule (skip confirmation)
  qshell sandbox transform-rule delete rule-xxxxxxxxxxxx -y
  qshell sbx tr dl rule-xxxxxxxxxxxx -y

  # Delete multiple rules
  qshell sandbox transform-rule delete rule-aaa rule-bbb -y
  qshell sbx tr dl rule-aaa rule-bbb -y

  # Interactively select rules to delete
  qshell sandbox transform-rule delete -s
  qshell sbx tr dl -s`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTransformRuleDeleteType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) == 0 && !info.Select {
				_ = cmd.Usage()
				return
			}
			info.RuleIDs = args
			operations.Delete(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Yes, "yes", "y", false, "skip confirmation")
	cmd.Flags().BoolVarP(&info.Select, "select", "s", false, "interactively select rules to delete")
	return cmd
}

// transformRuleCmdLoader adds the transform-rule command and its subcommands to the given parent command.
func transformRuleCmdLoader(parentCmd *cobra.Command, cfg *iqshell.Config) {
	transformRuleCmd := transformRuleCmdBuilder(cfg)
	transformRuleCmd.AddCommand(
		transformRuleListCmdBuilder(cfg),
		transformRuleGetCmdBuilder(cfg),
		transformRuleCreateCmdBuilder(cfg),
		transformRuleUpdateCmdBuilder(cfg),
		transformRuleDeleteCmdBuilder(cfg),
	)
	parentCmd.AddCommand(transformRuleCmd)
}
