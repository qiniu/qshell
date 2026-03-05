package operations

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// PublishInfo holds parameters for publishing/unpublishing templates.
type PublishInfo struct {
	TemplateIDs []string // One or more template IDs
	Yes         bool     // Skip confirmation
	Select      bool     // Interactive multi-select from template list
	Public      bool     // true = publish, false = unpublish
}

// Publish publishes or unpublishes one or more templates.
func Publish(info PublishInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	templateIDs := info.TemplateIDs

	action := "publish"
	if !info.Public {
		action = "unpublish"
	}

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
			publicStr := "private"
			if t.Public {
				publicStr = "public"
			}
			label = fmt.Sprintf("%s [%s]", label, publicStr)
			options = append(options, huh.NewOption(label, t.TemplateID))
		}

		var selected []string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title(fmt.Sprintf("Select templates to %s", action)).
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
		fmt.Printf("Are you sure you want to %s %d template(s)? [y/N] ", action, len(templateIDs))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Aborted")
			return
		}
	}

	for _, id := range templateIDs {
		if uErr := client.UpdateTemplate(ctx, id, sandbox.UpdateTemplateParams{
			Public: &info.Public,
		}); uErr != nil {
			sbClient.PrintError("%s template %s failed: %v", action, id, uErr)
			continue
		}
		if info.Public {
			sbClient.PrintSuccess("Template %s published", id)
		} else {
			sbClient.PrintSuccess("Template %s unpublished", id)
		}
	}
}
