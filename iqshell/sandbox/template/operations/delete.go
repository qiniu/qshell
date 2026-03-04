package operations

import (
	"context"
	"fmt"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// DeleteInfo holds parameters for deleting a template.
type DeleteInfo struct {
	TemplateID string
	Yes        bool // Skip confirmation
}

// Delete deletes a template.
func Delete(info DeleteInfo) {
	if info.TemplateID == "" {
		fmt.Println("Error: template ID is required")
		return
	}

	if !info.Yes {
		fmt.Printf("Are you sure you want to delete template %s? [y/N] ", info.TemplateID)
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Aborted")
			return
		}
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := client.DeleteTemplate(context.Background(), info.TemplateID); err != nil {
		fmt.Printf("Error: delete template failed: %v\n", err)
		return
	}

	fmt.Printf("Template %s deleted\n", info.TemplateID)
}
