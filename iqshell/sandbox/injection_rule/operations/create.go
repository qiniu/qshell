package operations

import (
	"context"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// CreateInfo holds parameters for creating an injection rule.
type CreateInfo struct {
	Name    string
	Type    string
	APIKey  string
	BaseURL string
	Headers string
}

// Create creates a new injection rule.
func Create(info CreateInfo) {
	if info.Name == "" {
		sbClient.PrintError("--name is required")
		return
	}

	spec, err := buildInjectionSpec(injectionInput{
		Type:    info.Type,
		APIKey:  info.APIKey,
		BaseURL: info.BaseURL,
		Headers: info.Headers,
	})
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	rule, err := client.CreateInjectionRule(context.Background(), sandbox.CreateInjectionRuleParams{
		Name:      info.Name,
		Injection: spec,
	})
	if err != nil {
		sbClient.PrintError("create injection rule failed: %v", err)
		return
	}

	sbClient.PrintSuccess("Injection rule created: %s (%s)", rule.RuleID, rule.Name)
}
