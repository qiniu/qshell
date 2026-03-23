package operations

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// GetInfo holds parameters for getting transform rule details.
type GetInfo struct {
	RuleID string
}

// Get retrieves and displays transform rule details.
func Get(info GetInfo) {
	if info.RuleID == "" {
		sbClient.PrintError("rule ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	rule, err := client.GetTransformRule(context.Background(), info.RuleID)
	if err != nil {
		sbClient.PrintError("get transform rule failed: %v", err)
		return
	}

	fmt.Printf("Rule ID:      %s\n", rule.RuleID)
	fmt.Printf("Name:         %s\n", rule.Name)
	fmt.Printf("Created At:   %s\n", rule.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated At:   %s\n", rule.UpdatedAt.Format(time.RFC3339))

	if rule.Conditions != nil && rule.Conditions.Hosts != nil && len(*rule.Conditions.Hosts) > 0 {
		fmt.Printf("\nConditions:\n")
		fmt.Printf("  Hosts:      %s\n", strings.Join(*rule.Conditions.Hosts, ", "))
	}

	printReplacements(rule.Replacements)
}

// printReplacements 输出替换动作的详情。
func printReplacements(r *sandbox.RequestTransformReplacements) {
	if r == nil {
		return
	}

	hasHeaders := r.Headers != nil && len(*r.Headers) > 0
	hasQueries := r.Queries != nil && len(*r.Queries) > 0
	if !hasHeaders && !hasQueries {
		return
	}

	fmt.Printf("\nReplacements:\n")
	if hasHeaders {
		fmt.Printf("  Headers:\n")
		for k, v := range *r.Headers {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
	if hasQueries {
		fmt.Printf("  Queries:\n")
		for k, v := range *r.Queries {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
}
