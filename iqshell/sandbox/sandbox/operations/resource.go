package operations

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ResourceListInfo holds parameters for listing resources mounted in a sandbox.
type ResourceListInfo struct {
	SandboxID string
	Format    string
}

// ResourceUpdateInfo holds parameters for updating a Git repository resource token.
type ResourceUpdateInfo struct {
	SandboxID  string
	ResourceID string
	Token      string
}

type resourceRow struct {
	ID        string
	Type      string
	Target    string
	MountPath string
	ReadOnly  string
}

// ListResources lists resources mounted in a sandbox.
func ListResources(info ResourceListInfo) {
	if strings.TrimSpace(info.SandboxID) == "" {
		sbClient.PrintError("sandbox ID is required")
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	resources, err := client.GetResources(ctx, info.SandboxID)
	if err != nil {
		sbClient.PrintError("list resources for sandbox %s failed: %v", info.SandboxID, err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(resources)
		return
	}

	if len(resources) == 0 {
		fmt.Println("No resources found")
		return
	}

	tw := sbClient.NewTable(os.Stdout)
	fmt.Fprintln(tw, "RESOURCE ID\tTYPE\tTARGET\tMOUNT PATH\tREAD ONLY")
	for _, resource := range resources {
		row, rowErr := resourceRowFromInfo(resource)
		if rowErr != nil {
			sbClient.PrintError("format resource for sandbox %s failed: %v", info.SandboxID, rowErr)
			return
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			row.ID,
			row.Type,
			row.Target,
			row.MountPath,
			row.ReadOnly,
		)
	}
	tw.Flush()
}

// UpdateResourceToken updates the authorization token of a Git repository resource.
func UpdateResourceToken(info ResourceUpdateInfo) {
	if err := validateResourceUpdateInfo(info); err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	ctx := context.Background()
	if err := client.UpdateGitRepositoryResourceToken(ctx, info.SandboxID, info.ResourceID, info.Token); err != nil {
		sbClient.PrintError("update resource %s token failed: %v", info.ResourceID, err)
		return
	}
	sbClient.PrintSuccess("Updated Git repository resource %s token", info.ResourceID)
}

func validateResourceUpdateInfo(info ResourceUpdateInfo) error {
	switch {
	case strings.TrimSpace(info.SandboxID) == "":
		return fmt.Errorf("sandbox ID is required")
	case strings.TrimSpace(info.ResourceID) == "":
		return fmt.Errorf("resource ID is required")
	case strings.TrimSpace(info.Token) == "":
		return fmt.Errorf("token is required")
	default:
		return nil
	}
}

func resourceRowFromInfo(resource sandbox.SandboxResourceInfo) (resourceRow, error) {
	switch {
	case resource.GitRepository != nil:
		return resourceRow{
			ID:        resource.GitRepository.ResourceID,
			Type:      string(resource.GitRepository.Type),
			Target:    resource.GitRepository.URL,
			MountPath: resource.GitRepository.MountPath,
			ReadOnly:  "-",
		}, nil
	case resource.Kodo != nil:
		target := resource.Kodo.Bucket
		if resource.Kodo.Prefix != nil {
			target += "/" + strings.TrimPrefix(*resource.Kodo.Prefix, "/")
		}
		readOnly := "-"
		if resource.Kodo.ReadOnly != nil {
			readOnly = strconv.FormatBool(*resource.Kodo.ReadOnly)
		}
		return resourceRow{
			ID:        resource.Kodo.ResourceID,
			Type:      "kodo",
			Target:    target,
			MountPath: resource.Kodo.MountPath,
			ReadOnly:  readOnly,
		}, nil
	default:
		return resourceRow{}, fmt.Errorf("unknown sandbox resource type")
	}
}
