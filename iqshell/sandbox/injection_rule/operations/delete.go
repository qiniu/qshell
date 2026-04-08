package operations

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// DeleteInfo holds parameters for deleting injection rules.
type DeleteInfo struct {
	RuleIDs []string
	Yes     bool
	Select  bool
}

// Delete deletes one or more injection rules.
func Delete(info DeleteInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	ruleIDs := info.RuleIDs

	if info.Select {
		rules, lErr := client.ListInjectionRules(ctx)
		if lErr != nil {
			sbClient.PrintError("list injection rules failed: %v", lErr)
			return
		}
		if len(rules) == 0 {
			fmt.Println("No injection rules found")
			return
		}

		options := make([]huh.Option[string], 0, len(rules))
		for _, r := range rules {
			label := fmt.Sprintf("%s (%s)", r.RuleID, r.Name)
			options = append(options, huh.NewOption(label, r.RuleID))
		}

		var selected []string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select injection rules to delete").
					Options(options...).
					Value(&selected),
			),
		)
		if fErr := form.Run(); fErr != nil {
			sbClient.PrintError("selection cancelled: %v", fErr)
			return
		}
		if len(selected) == 0 {
			fmt.Println("No injection rules selected")
			return
		}
		ruleIDs = selected
	}

	if len(ruleIDs) == 0 {
		sbClient.PrintError("at least one rule ID is required (or use --select)")
		return
	}

	if !info.Yes {
		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete %d injection rule(s)?", len(ruleIDs))).
					Value(&confirm),
			),
		)
		if fErr := form.Run(); fErr != nil || !confirm {
			fmt.Println("Aborted")
			return
		}
	}

	hasError := false
	for _, id := range ruleIDs {
		if dErr := client.DeleteInjectionRule(ctx, id); dErr != nil {
			sbClient.PrintError("delete injection rule %s failed: %v", id, dErr)
			hasError = true
			continue
		}
		sbClient.PrintSuccess("Injection rule %s deleted", id)
	}
	if hasError {
		os.Exit(1)
	}
}
