package operations

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// DeleteInfo holds parameters for deleting transform rules.
type DeleteInfo struct {
	RuleIDs []string // 一个或多个规则 ID
	Yes     bool     // 跳过确认
	Select  bool     // 交互式多选
}

// Delete deletes one or more transform rules.
func Delete(info DeleteInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	ruleIDs := info.RuleIDs

	// 交互式选择模式
	if info.Select {
		rules, lErr := client.ListTransformRules(ctx)
		if lErr != nil {
			sbClient.PrintError("list transform rules failed: %v", lErr)
			return
		}
		if len(rules) == 0 {
			fmt.Println("No transform rules found")
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
					Title("Select transform rules to delete").
					Options(options...).
					Value(&selected),
			),
		)
		if fErr := form.Run(); fErr != nil {
			sbClient.PrintError("selection cancelled: %v", fErr)
			return
		}
		if len(selected) == 0 {
			fmt.Println("No transform rules selected")
			return
		}
		ruleIDs = selected
	}

	if len(ruleIDs) == 0 {
		sbClient.PrintError("at least one rule ID is required (or use --select)")
		return
	}

	if !info.Yes {
		fmt.Printf("Are you sure you want to delete %d transform rule(s)? [y/N] ", len(ruleIDs))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Aborted")
			return
		}
	}

	for _, id := range ruleIDs {
		if dErr := client.DeleteTransformRule(ctx, id); dErr != nil {
			sbClient.PrintError("delete transform rule %s failed: %v", id, dErr)
			continue
		}
		sbClient.PrintSuccess("Transform rule %s deleted", id)
	}
}
