package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMinDiskPreservesExplicitZeroAndCLIOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "qshell.sandbox.toml")
	require.NoError(t, os.WriteFile(path, []byte("min_disk_mb = 5120\n"), 0600))
	cfg, err := Load(path)
	require.NoError(t, err)
	fields := BuildFields{}
	cfg.ApplyTo(&fields)
	require.NotNil(t, fields.MinDiskMB)
	require.EqualValues(t, 5120, *fields.MinDiskMB)
	zero := int32(0)
	fields = BuildFields{MinDiskMB: &zero}
	require.Contains(t, cfg.ApplyTo(&fields), "min_disk_mb")
	require.Zero(t, *fields.MinDiskMB)
	require.NoError(t, os.WriteFile(path, []byte("min_disk_mb = 0\n"), 0600))
	cfg, err = Load(path)
	require.NoError(t, err)
	fields = BuildFields{}
	cfg.ApplyTo(&fields)
	require.NotNil(t, fields.MinDiskMB)
	require.Zero(t, *fields.MinDiskMB)
}
