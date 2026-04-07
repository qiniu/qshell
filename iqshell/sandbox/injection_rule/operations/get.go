package operations

import (
	"context"
	"fmt"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// GetInfo holds parameters for getting injection rule details.
type GetInfo struct {
	RuleID string
}

// Get retrieves and displays injection rule details.
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

	rule, err := client.GetInjectionRule(context.Background(), info.RuleID)
	if err != nil {
		sbClient.PrintError("get injection rule failed: %v", err)
		return
	}

	fmt.Printf("Rule ID:      %s\n", rule.RuleID)
	fmt.Printf("Name:         %s\n", rule.Name)
	fmt.Printf("Type:         %s\n", formatInjectionType(rule.Injection))
	fmt.Printf("Target:       %s\n", formatInjectionTarget(rule.Injection))
	fmt.Printf("API Key:      %s\n", yesNo(hasAPIKey(rule.Injection)))
	fmt.Printf("Headers:      %s\n", formatInjectionHeaders(rule.Injection))
	fmt.Printf("Created At:   %s\n", sbClient.FormatTimestamp(rule.CreatedAt))
	fmt.Printf("Updated At:   %s\n", sbClient.FormatTimestamp(rule.UpdatedAt))
}

func yesNo(value bool) string {
	if value {
		return "configured"
	}
	return "-"
}
