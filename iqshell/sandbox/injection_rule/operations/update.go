package operations

import (
	"context"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// UpdateInfo holds parameters for updating an injection rule.
type UpdateInfo struct {
	RuleID  string
	Name    string
	Type    string
	APIKey  string
	BaseURL string
	Headers string
}

// Update updates an existing injection rule.
func Update(info UpdateInfo) {
	if info.RuleID == "" {
		sbClient.PrintError("rule ID is required")
		return
	}

	input := injectionInput{
		Type:    info.Type,
		APIKey:  info.APIKey,
		BaseURL: info.BaseURL,
		Headers: info.Headers,
	}

	if info.Name == "" && !shouldUpdateInjection(input) {
		sbClient.PrintError("at least one of --name, --type, --api-key, --base-url, or --headers is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	params := sandbox.UpdateInjectionRuleParams{}
	if info.Name != "" {
		params.Name = &info.Name
	}
	if shouldUpdateInjection(input) {
		spec, err := buildInjectionSpec(input)
		if err != nil {
			sbClient.PrintError("%v", err)
			return
		}
		params.Injection = &spec
	}

	rule, err := client.UpdateInjectionRule(context.Background(), info.RuleID, params)
	if err != nil {
		sbClient.PrintError("update injection rule failed: %v", err)
		return
	}

	sbClient.PrintSuccess("Injection rule updated: %s (%s)", rule.RuleID, rule.Name)
}
