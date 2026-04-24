package operations

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// DeleteInfo holds parameters for deleting templates.
type DeleteInfo struct {
	TemplateIDs []string // One or more template IDs to delete
	Yes         bool     // Skip confirmation
	Select      bool     // Interactive multi-select from template list
}

// Delete deletes one or more templates.
func Delete(info DeleteInfo) {
	if len(info.TemplateIDs) == 0 && !info.Select {
		id, ok := templateIDFromCwdConfig()
		if !ok {
			return
		}
		if id != "" {
			info.TemplateIDs = []string{id}
		}
	}
	if len(info.TemplateIDs) == 0 && !info.Select {
		sbClient.PrintError("at least one template ID is required (positional args, --select, or qshell.sandbox.toml)")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	templateIDs := info.TemplateIDs

	// Interactive selection mode
	if info.Select {
		templates, lErr := client.ListTemplates(ctx, nil)
		if lErr != nil {
			sbClient.PrintError("list templates failed: %v", lErr)
			return
		}
		if len(templates) == 0 {
			fmt.Println("No templates found")
			return
		}

		options := make([]huh.Option[string], 0, len(templates))
		for _, t := range templates {
			label := t.TemplateID
			if len(t.Aliases) > 0 {
				label = fmt.Sprintf("%s (%s)", t.TemplateID, t.Aliases[0])
			}
			options = append(options, huh.NewOption(label, t.TemplateID))
		}

		var selected []string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select templates to delete").
					Options(options...).
					Value(&selected),
			),
		)
		if fErr := form.Run(); fErr != nil {
			sbClient.PrintError("selection cancelled: %v", fErr)
			return
		}
		if len(selected) == 0 {
			fmt.Println("No templates selected")
			return
		}
		templateIDs = selected
	}

	if len(templateIDs) == 0 {
		sbClient.PrintError("at least one template ID is required (or use --select)")
		return
	}

	if !info.Yes {
		fmt.Printf("Are you sure you want to delete %d template(s)? [y/N] ", len(templateIDs))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Aborted")
			return
		}
	}

	for _, id := range templateIDs {
		if dErr := client.DeleteTemplate(ctx, id); dErr != nil {
			sbClient.PrintError("delete template %s failed: %v", id, dErr)
			continue
		}
		sbClient.PrintSuccess("Template %s deleted", id)
	}
}
