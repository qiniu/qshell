package cmd

import "github.com/spf13/cobra"

type cmdBuilder func() *cobra.Command
