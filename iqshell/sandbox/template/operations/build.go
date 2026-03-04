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
}

// Build creates or rebuilds a template.
// If TemplateID is provided, starts a new build for the existing template.
// Otherwise, creates a new template with the given Name and starts a build.
func Build(info BuildInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ctx := context.Background()
	templateID := info.TemplateID
	buildID := ""

	if templateID == "" {
		// Create a new template
		if info.Name == "" {
			fmt.Println("Error: template name (--name) or template ID (--template-id) is required")
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
			fmt.Printf("Error: create template failed: %v\n", cErr)
			return
		}
		templateID = resp.TemplateID
		buildID = resp.BuildID
		fmt.Printf("Template %s created (build ID: %s)\n", templateID, buildID)
	} else {
		// Get existing template to find latest build ID
		tmpl, gErr := client.GetTemplate(ctx, templateID, nil)
		if gErr != nil {
			fmt.Printf("Error: get template failed: %v\n", gErr)
			return
		}
		if len(tmpl.Builds) > 0 {
			buildID = tmpl.Builds[0].BuildID
		} else {
			fmt.Println("Error: no builds found for template, cannot rebuild")
			return
		}
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

	fmt.Printf("Starting build for template %s (build ID: %s)...\n", templateID, buildID)
	if err := client.StartTemplateBuild(ctx, templateID, buildID, buildParams); err != nil {
		fmt.Printf("Error: start build failed: %v\n", err)
		return
	}

	if !info.Wait {
		fmt.Printf("Build started. Use 'qshell sandbox template builds %s %s' to check status.\n", templateID, buildID)
		return
	}

	// Wait for build completion
	fmt.Println("Waiting for build to complete...")
	buildInfo, err := client.WaitForBuild(ctx, templateID, buildID,
		sandbox.WithPollInterval(3*time.Second),
	)
	if err != nil {
		fmt.Printf("Error: build failed: %v\n", err)
		return
	}

	fmt.Printf("Build completed!\n")
	fmt.Printf("Template ID:  %s\n", buildInfo.TemplateID)
	fmt.Printf("Build ID:     %s\n", buildInfo.BuildID)
	fmt.Printf("Status:       %s\n", buildInfo.Status)
}
