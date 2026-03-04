package operations

import (
	"context"
	"fmt"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// BuildsInfo holds parameters for viewing template build status.
type BuildsInfo struct {
	TemplateID string
	BuildID    string
}

// Builds retrieves and displays the build status of a template.
func Builds(info BuildsInfo) {
	if info.TemplateID == "" {
		fmt.Println("Error: template ID is required")
		return
	}
	if info.BuildID == "" {
		fmt.Println("Error: build ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	buildInfo, err := client.GetTemplateBuildStatus(context.Background(), info.TemplateID, info.BuildID, nil)
	if err != nil {
		fmt.Printf("Error: get build status failed: %v\n", err)
		return
	}

	fmt.Printf("Template ID:  %s\n", buildInfo.TemplateID)
	fmt.Printf("Build ID:     %s\n", buildInfo.BuildID)
	fmt.Printf("Status:       %s\n", buildInfo.Status)

	if len(buildInfo.Logs) > 0 {
		fmt.Printf("\nBuild Logs:\n")
		for _, log := range buildInfo.Logs {
			fmt.Printf("  %s\n", log)
		}
	}
}
