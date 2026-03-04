package operations

import (
	"context"
	"fmt"
	"time"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ListInfo holds parameters for listing templates.
type ListInfo struct {
	Format string // pretty or json
}

// List lists all templates.
func List(info ListInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	templates, err := client.ListTemplates(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error: list templates failed: %v\n", err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(templates)
		return
	}

	if len(templates) == 0 {
		fmt.Println("No templates found")
		return
	}

	fmt.Printf("%-30s %-20s %-10s %-6s %-10s %-10s %s\n",
		"TEMPLATE ID", "ALIASES", "STATUS", "CPU", "MEMORY", "DISK", "UPDATED AT")
	for _, t := range templates {
		aliases := "-"
		if len(t.Aliases) > 0 {
			aliases = t.Aliases[0]
		}
		fmt.Printf("%-30s %-20s %-10s %-6d %-10s %-10s %s\n",
			t.TemplateID,
			aliases,
			t.BuildStatus,
			t.CPUCount,
			fmt.Sprintf("%dMB", t.MemoryMB),
			fmt.Sprintf("%dMB", t.DiskSizeMB),
			t.UpdatedAt.Format(time.RFC3339),
		)
	}
}
