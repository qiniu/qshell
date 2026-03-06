package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// BuildInfo holds parameters for building a template.
type BuildInfo struct {
	// Name is the template name (used for create + build).
	Name string

	// TemplateID is the existing template ID (used for rebuild).
	TemplateID string

	// FromImage is the base Docker image.
	FromImage string

	// FromTemplate is the base template.
	FromTemplate string

	// StartCmd is the command to run after build.
	StartCmd string

	// ReadyCmd is the readiness check command.
	ReadyCmd string

	// CPUCount is the sandbox CPU count.
	CPUCount int32

	// MemoryMB is the sandbox memory size in MiB.
	MemoryMB int32

	// Wait indicates whether to wait for build completion.
	Wait bool

	// NoCache forces a full rebuild ignoring cache.
	NoCache bool
}

// Build creates or rebuilds a template.
// If TemplateID is provided, starts a new build for the existing template.
// Otherwise, creates a new template with the given Name and starts a build.
func Build(info BuildInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	templateID := info.TemplateID
	buildID := ""

	if templateID == "" {
		// Create a new template
		if info.Name == "" {
			sbClient.PrintError("template name (--name) or template ID (--template-id) is required")
			return
		}

		createParams := sandbox.CreateTemplateParams{
			Name: &info.Name,
		}
		if info.CPUCount > 0 {
			createParams.CPUCount = &info.CPUCount
		}
		if info.MemoryMB > 0 {
			createParams.MemoryMB = &info.MemoryMB
		}

		fmt.Printf("Creating template %s...\n", info.Name)
		resp, cErr := client.CreateTemplate(ctx, createParams)
		if cErr != nil {
			sbClient.PrintError("create template failed: %v", cErr)
			return
		}
		templateID = resp.TemplateID
		buildID = resp.BuildID
		sbClient.PrintSuccess("Template %s created (build ID: %s)", templateID, buildID)
	} else {
		// Get existing template to find latest build ID
		tmpl, gErr := client.GetTemplate(ctx, templateID, nil)
		if gErr != nil {
			sbClient.PrintError("get template failed: %v", gErr)
			return
		}
		if len(tmpl.Builds) > 0 {
			buildID = tmpl.Builds[0].BuildID
		} else {
			sbClient.PrintError("no builds found for template, cannot rebuild")
			return
		}
	}

	// Validate source
	if info.FromImage == "" && info.FromTemplate == "" {
		sbClient.PrintError("--from-image or --from-template is required")
		return
	}

	// Start build
	buildParams := sandbox.StartTemplateBuildParams{}
	if info.FromImage != "" {
		buildParams.FromImage = &info.FromImage
	}
	if info.FromTemplate != "" {
		buildParams.FromTemplate = &info.FromTemplate
	}
	if info.StartCmd != "" {
		buildParams.StartCmd = &info.StartCmd
	}
	if info.ReadyCmd != "" {
		buildParams.ReadyCmd = &info.ReadyCmd
	}
	if info.NoCache {
		force := true
		buildParams.Force = &force
	}

	fmt.Printf("Starting build for template %s (build ID: %s)...\n", templateID, buildID)
	if err := client.StartTemplateBuild(ctx, templateID, buildID, buildParams); err != nil {
		sbClient.PrintError("start build failed: %v", err)
		return
	}

	if !info.Wait {
		fmt.Printf("Build started. Use 'qshell sandbox template builds %s %s' to check status.\n", templateID, buildID)
		return
	}

	// Stream build logs while waiting
	fmt.Println("Waiting for build to complete...")
	var cursor *int64
	for {
		logs, blErr := client.GetTemplateBuildLogs(ctx, templateID, buildID, &sandbox.GetBuildLogsParams{
			Cursor: cursor,
		})
		if blErr == nil && logs != nil {
			for _, entry := range logs.Logs {
				fmt.Printf("[%s] %s %s\n",
					sbClient.FormatTimestamp(entry.Timestamp),
					sbClient.LogLevelBadge(string(entry.Level)),
					entry.Message,
				)
				ts := entry.Timestamp.UnixMilli() + 1
				cursor = &ts
			}
		}

		// Check build status
		buildInfo, bErr := client.GetTemplateBuildStatus(ctx, templateID, buildID, nil)
		if bErr != nil {
			sbClient.PrintError("get build status failed: %v", bErr)
			return
		}

		if buildInfo.Status == "ready" || buildInfo.Status == "error" {
			if buildInfo.Status == "error" {
				sbClient.PrintError("build failed")
			} else {
				sbClient.PrintSuccess("Build completed!")
			}
			fmt.Printf("Template ID:  %s\n", buildInfo.TemplateID)
			fmt.Printf("Build ID:     %s\n", buildInfo.BuildID)
			fmt.Printf("Status:       %s\n", buildInfo.Status)
			return
		}

		time.Sleep(3 * time.Second)
	}
}
