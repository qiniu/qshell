package operations

import (
	"context"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// UpdateInfo holds parameters for updating a transform rule.
type UpdateInfo struct {
	RuleID  string // 规则 ID（必填）
	Name    string // 规则名称
	Hosts   string // 逗号分隔的域名列表
	Headers string // 逗号分隔的 key=value 对
	Queries string // 逗号分隔的 key=value 对
}

// Update updates an existing transform rule.
func Update(info UpdateInfo) {
	if info.RuleID == "" {
		sbClient.PrintError("rule ID is required")
		return
	}

	if info.Name == "" && info.Hosts == "" && info.Headers == "" && info.Queries == "" {
		sbClient.PrintError("at least one of --name, --hosts, --headers, or --queries is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	params := sandbox.UpdateTransformRuleParams{}

	if info.Name != "" {
		params.Name = &info.Name
	}

	if info.Hosts != "" {
		hosts := parseCommaSeparated(info.Hosts)
		params.Conditions = &sandbox.RequestTransformConditions{
			Hosts: &hosts,
		}
	}

	if info.Headers != "" || info.Queries != "" {
		replacements := &sandbox.RequestTransformReplacements{}
		if info.Headers != "" {
			headers := sbClient.ParseMetadataMap(info.Headers)
			replacements.Headers = &headers
		}
		if info.Queries != "" {
			queries := sbClient.ParseMetadataMap(info.Queries)
			replacements.Queries = &queries
		}
		params.Replacements = replacements
	}

	rule, err := client.UpdateTransformRule(context.Background(), info.RuleID, params)
	if err != nil {
		sbClient.PrintError("update transform rule failed: %v", err)
		return
	}

	sbClient.PrintSuccess("Transform rule updated: %s (%s)", rule.RuleID, rule.Name)
}
