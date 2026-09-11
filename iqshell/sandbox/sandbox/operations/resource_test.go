//go:build unit

package operations

import (
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

func TestResourceRowFromInfo_GitRepository(t *testing.T) {
	row, err := resourceRowFromInfo(sandbox.SandboxResourceInfo{
		GitRepository: &sandbox.GitRepositoryResourceInfo{
			ResourceID: "res_git",
			Type:       sandbox.GitRepositoryTypeGithub,
			URL:        "https://github.com/qiniu/qshell.git",
			MountPath:  "/workspace/qshell",
		},
	})
	if err != nil {
		t.Fatalf("resourceRowFromInfo() error = %v", err)
	}
	if row.ID != "res_git" || row.Type != "github_repository" {
		t.Fatalf("row identity = %+v", row)
	}
	if row.Target != "https://github.com/qiniu/qshell.git" || row.MountPath != "/workspace/qshell" {
		t.Fatalf("row location = %+v", row)
	}
	if row.ReadOnly != "-" {
		t.Fatalf("git read-only = %q, want -", row.ReadOnly)
	}
}

func TestResourceRowFromInfo_Kodo(t *testing.T) {
	prefix := "datasets/"
	readOnly := true
	row, err := resourceRowFromInfo(sandbox.SandboxResourceInfo{
		Kodo: &sandbox.KodoResourceInfo{
			ResourceID: "res_kodo",
			Bucket:     "data-bucket",
			MountPath:  "/mnt/data",
			Prefix:     &prefix,
			ReadOnly:   &readOnly,
		},
	})
	if err != nil {
		t.Fatalf("resourceRowFromInfo() error = %v", err)
	}
	if row.ID != "res_kodo" || row.Type != "kodo" {
		t.Fatalf("row identity = %+v", row)
	}
	if row.Target != "data-bucket/datasets/" || row.MountPath != "/mnt/data" {
		t.Fatalf("row location = %+v", row)
	}
	if row.ReadOnly != "true" {
		t.Fatalf("kodo read-only = %q, want true", row.ReadOnly)
	}
}

func TestResourceRowFromInfo_RejectsUnknownResource(t *testing.T) {
	if _, err := resourceRowFromInfo(sandbox.SandboxResourceInfo{}); err == nil {
		t.Fatal("expected unknown resource to fail")
	}
}

func TestValidateResourceUpdateInfo(t *testing.T) {
	tests := []struct {
		name string
		info ResourceUpdateInfo
	}{
		{name: "missing sandbox ID", info: ResourceUpdateInfo{ResourceID: "res", Token: "token"}},
		{name: "missing resource ID", info: ResourceUpdateInfo{SandboxID: "sb", Token: "token"}},
		{name: "missing token", info: ResourceUpdateInfo{SandboxID: "sb", ResourceID: "res"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateResourceUpdateInfo(tt.info); err == nil {
				t.Fatal("expected invalid update info to fail")
			}
		})
	}
}
