package operations

import (
	"context"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// CreateInfo holds parameters for creating a transform rule.
type CreateInfo struct {
	Name    string // 规则名称（必填）
	Hosts   string // 逗号分隔的域名列表
	Headers string // 逗号分隔的 key=value 对
	Queries string // 逗号分隔的 key=value 对
}

// Create creates a new transform rule.
func Create(info CreateInfo) {
	if info.Name == "" {
		sbClient.PrintError("--name is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	params := sandbox.CreateTransformRuleParams{
		Name: info.Name,
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

	rule, err := client.CreateTransformRule(context.Background(), params)
	if err != nil {
		sbClient.PrintError("create transform rule failed: %v", err)
		return
	}

	sbClient.PrintSuccess("Transform rule created: %s (%s)", rule.RuleID, rule.Name)
}

// parseCommaSeparated 将逗号分隔的字符串解析为字符串切片。
func parseCommaSeparated(raw string) []string {
	var result []string
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
