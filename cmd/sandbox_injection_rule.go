package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/injection_rule/operations"
)

var injectionRuleCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "injection-rule",
		Aliases: []string{"ir"},
		Short:   "Manage sandbox injection rules (alias: ir)",
		Args:    cobra.NoArgs,
		Example: `  # View injection-rule subcommands
  qshell sandbox injection-rule -h
  qshell sbx ir -h

  # List all injection rules
  qshell sandbox injection-rule list
  qshell sbx ir ls

  # Create an OpenAI injection rule
  qshell sandbox injection-rule create --name openai-default --type openai --api-key sk-xxx
  qshell sbx ir cr --name openai-default --type openai --api-key sk-xxx

  # Get injection rule details
  qshell sandbox injection-rule get rule-xxxxxxxxxxxx
  qshell sbx ir gt rule-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxInjectionRuleType
			docs.ShowCmdDocument(docs.SandboxInjectionRuleType)
		},
	}
	return cmd
}

var injectionRuleListCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ListInfo{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List injection rules (alias: ls)",
		Example: `  # List all injection rules
  qshell sandbox injection-rule list
  qshell sbx ir ls

  # Output as JSON
  qshell sandbox injection-rule list --format json
  qshell sbx ir ls --format json`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxInjectionRuleListType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.List(info)
		},
	}
	cmd.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	return cmd
}

var injectionRuleGetCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <ruleID>",
		Aliases: []string{"gt"},
		Short:   "Get injection rule details (alias: gt)",
		Args:    cobra.ExactArgs(1),
		Example: `  # Get injection rule details
  qshell sandbox injection-rule get rule-xxxxxxxxxxxx
  qshell sbx ir gt rule-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxInjectionRuleGetType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.Get(operations.GetInfo{
				RuleID: args[0],
			})
		},
	}
	return cmd
}

var injectionRuleCreateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.CreateInfo{}
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"cr"},
		Short:   "Create an injection rule (alias: cr)",
		Example: `  # Create an OpenAI injection rule
  qshell sandbox injection-rule create --name openai-default --type openai --api-key sk-xxx
  qshell sbx ir cr --name openai-default --type openai --api-key sk-xxx

  # Create an Anthropic injection rule with custom base URL
  qshell sandbox injection-rule create --name anthropic-proxy --type anthropic --api-key sk-ant --base-url https://anthropic-proxy.example.com

  # Create a custom HTTP injection rule
  qshell sandbox injection-rule create --name api-auth --type http --base-url https://api.example.com --headers "Authorization=Bearer token123,X-Env=prod"
  qshell sbx ir cr --name api-auth --type http --base-url https://api.example.com --headers "Authorization=Bearer token123,X-Env=prod"

  # Create a Qiniu AI API injection rule
  qshell sandbox injection-rule create --name qiniu-ai --type qiniu --api-key ak-xxx
  qshell sbx ir cr --name qiniu-ai --type qiniu --api-key ak-xxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxInjectionRuleCreateType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.Create(info)
		},
	}
	cmd.Flags().StringVar(&info.Name, "name", "", "rule name (required, unique per user)")
	cmd.Flags().StringVar(&info.Type, "type", "", "injection type: openai, anthropic, gemini, qiniu, http")
	cmd.Flags().StringVar(&info.APIKey, "api-key", "", "API key for openai/anthropic/gemini/qiniu injection types (warning: passing secrets via CLI may leak through shell history or process lists)")
	cmd.Flags().StringVar(&info.BaseURL, "base-url", "", "override base URL or target base URL for http injection")
	cmd.Flags().StringVar(&info.Headers, "headers", "", "HTTP headers for custom http injection (comma-separated key=value pairs)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

var injectionRuleUpdateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UpdateInfo{}
	cmd := &cobra.Command{
		Use:     "update <ruleID>",
		Aliases: []string{"up"},
		Short:   "Update an injection rule (alias: up)",
		Args:    cobra.ExactArgs(1),
		Example: `  # Update rule name
  qshell sandbox injection-rule update rule-xxxxxxxxxxxx --name new-name
  qshell sbx ir up rule-xxxxxxxxxxxx --name new-name

  # Update to a Gemini injection with custom base URL
  qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type gemini --api-key sk-gem --base-url https://gemini-proxy.example.com
  qshell sbx ir up rule-xxxxxxxxxxxx --type gemini --api-key sk-gem --base-url https://gemini-proxy.example.com

  # Update custom HTTP headers
  qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type http --base-url https://api.example.com --headers "Authorization=Bearer newtoken"
  qshell sbx ir up rule-xxxxxxxxxxxx --type http --base-url https://api.example.com --headers "Authorization=Bearer newtoken"

  # Update to a Qiniu AI API injection
  qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type qiniu --api-key ak-new
  qshell sbx ir up rule-xxxxxxxxxxxx --type qiniu --api-key ak-new`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxInjectionRuleUpdateType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			info.RuleID = args[0]
			operations.Update(info)
		},
	}
	cmd.Flags().StringVar(&info.Name, "name", "", "new rule name")
	cmd.Flags().StringVar(&info.Type, "type", "", "new injection type: openai, anthropic, gemini, qiniu, http")
	cmd.Flags().StringVar(&info.APIKey, "api-key", "", "new API key for openai/anthropic/gemini/qiniu injection types (warning: passing secrets via CLI may leak through shell history or process lists)")
	cmd.Flags().StringVar(&info.BaseURL, "base-url", "", "new base URL or target base URL for http injection")
	cmd.Flags().StringVar(&info.Headers, "headers", "", "new HTTP headers for custom http injection (comma-separated key=value pairs)")
	return cmd
}

var injectionRuleDeleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DeleteInfo{}
	cmd := &cobra.Command{
		Use:     "delete [ruleIDs...]",
		Aliases: []string{"dl"},
		Short:   "Delete one or more injection rules (alias: dl)",
		Example: `  # Delete a single rule (skip confirmation)
  qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
  qshell sbx ir dl rule-xxxxxxxxxxxx -y

  # Delete multiple rules
  qshell sandbox injection-rule delete rule-aaa rule-bbb -y
  qshell sbx ir dl rule-aaa rule-bbb -y

  # Interactively select rules to delete
  qshell sandbox injection-rule delete -s
  qshell sbx ir dl -s`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxInjectionRuleDeleteType
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

// injectionRuleCmdLoader adds the injection-rule command and its subcommands to the given parent command.
func injectionRuleCmdLoader(parentCmd *cobra.Command, cfg *iqshell.Config) {
	injectionRuleCmd := injectionRuleCmdBuilder(cfg)
	injectionRuleCmd.AddCommand(
		injectionRuleListCmdBuilder(cfg),
		injectionRuleGetCmdBuilder(cfg),
		injectionRuleCreateCmdBuilder(cfg),
		injectionRuleUpdateCmdBuilder(cfg),
		injectionRuleDeleteCmdBuilder(cfg),
	)
	parentCmd.AddCommand(injectionRuleCmd)
}
