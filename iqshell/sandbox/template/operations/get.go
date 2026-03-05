package operations

import (
	"context"
	"fmt"
	"time"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// GetInfo holds parameters for getting template details.
type GetInfo struct {
	TemplateID string
}

// Get retrieves and displays template details.
func Get(info GetInfo) {
	if info.TemplateID == "" {
		sbClient.PrintError("template ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	tmpl, err := client.GetTemplate(context.Background(), info.TemplateID, nil)
	if err != nil {
		sbClient.PrintError("get template failed: %v", err)
		return
	}

	fmt.Printf("Template ID:    %s\n", tmpl.TemplateID)
	fmt.Printf("Aliases:        %v\n", tmpl.Aliases)
	fmt.Printf("Public:         %v\n", tmpl.Public)
	fmt.Printf("Spawn Count:    %d\n", tmpl.SpawnCount)
	fmt.Printf("Created At:     %s\n", tmpl.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated At:     %s\n", tmpl.UpdatedAt.Format(time.RFC3339))
	if tmpl.LastSpawnedAt != nil {
		fmt.Printf("Last Spawned:   %s\n", tmpl.LastSpawnedAt.Format(time.RFC3339))
	}

	if len(tmpl.Builds) > 0 {
		fmt.Printf("\nBuilds:\n")
		fmt.Printf("  %-36s %-10s %-6s %-10s %-10s %s\n",
			"BUILD ID", "STATUS", "CPU", "MEMORY", "DISK", "CREATED AT")
		for _, b := range tmpl.Builds {
			disk := "-"
			if b.DiskSizeMB != nil {
				disk = fmt.Sprintf("%dMB", *b.DiskSizeMB)
			}
			fmt.Printf("  %-36s %-10s %-6d %-10s %-10s %s\n",
				b.BuildID,
				b.Status,
				b.CPUCount,
				fmt.Sprintf("%dMB", b.MemoryMB),
				disk,
				b.CreatedAt.Format(time.RFC3339),
			)
		}
	}
}
