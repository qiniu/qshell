package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnapshotCommandRequiresStableKey(t *testing.T) {
	command := sandboxSnapshotCmdBuilder(nil)
	create, _, err := command.Find([]string{"create"})
	require.NoError(t, err)
	require.NoError(t, create.ParseFlags([]string{"sb-1", "--memory=false", "--idempotency-key=capture-1"}))
	memory, err := create.Flags().GetBool("memory")
	require.NoError(t, err)
	require.False(t, memory)
	require.Equal(t, []string{"true"}, create.Flags().Lookup("idempotency-key").Annotations["cobra_annotation_bash_completion_one_required_flag"])
	list, _, err := command.Find([]string{"list"})
	require.NoError(t, err)
	require.NotNil(t, list.Flags().Lookup("next-token"))
}
