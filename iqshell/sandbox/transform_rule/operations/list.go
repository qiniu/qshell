package operations

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ListInfo holds parameters for listing transform rules.
type ListInfo struct {
	Format string // pretty or json
}

// List lists all transform rules.
func List(info ListInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	rules, err := client.ListTransformRules(context.Background())
	if err != nil {
		sbClient.PrintError("list transform rules failed: %v", err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(rules)
		return
	}

	if len(rules) == 0 {
		fmt.Println("No transform rules found")
		return
	}

	tw := sbClient.NewTable(os.Stdout)
	fmt.Fprintf(tw, "RULE ID\tNAME\tHOSTS\tHEADERS\tQUERIES\tCREATED AT\tUPDATED AT\n")
	for _, r := range rules {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.RuleID,
			r.Name,
			formatHosts(r.Conditions),
			formatHeaderKeys(r.Replacements),
			formatQueryKeys(r.Replacements),
			sbClient.FormatTimestamp(r.CreatedAt),
			sbClient.FormatTimestamp(r.UpdatedAt),
		)
	}
	tw.Flush()
}

// formatHosts 格式化匹配条件中的域名列表。
func formatHosts(conditions *sandbox.RequestTransformConditions) string {
	if conditions == nil || conditions.Hosts == nil || len(*conditions.Hosts) == 0 {
		return "-"
	}
	return strings.Join(*conditions.Hosts, ", ")
}

// formatHeaderKeys 格式化替换动作中的 Headers 键列表。
func formatHeaderKeys(replacements *sandbox.RequestTransformReplacements) string {
	if replacements == nil || replacements.Headers == nil || len(*replacements.Headers) == 0 {
		return "-"
	}
	keys := make([]string, 0, len(*replacements.Headers))
	for k := range *replacements.Headers {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

// formatQueryKeys 格式化替换动作中的 Queries 键列表。
func formatQueryKeys(replacements *sandbox.RequestTransformReplacements) string {
	if replacements == nil || replacements.Queries == nil || len(*replacements.Queries) == 0 {
		return "-"
	}
	keys := make([]string, 0, len(*replacements.Queries))
	for k := range *replacements.Queries {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
