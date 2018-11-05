package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
	"os"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var userChCmd = &cobra.Command{
	Use:   "cu [<AccountName>]",
	Short: "Change user to AccountName",
	Args:  cobra.RangeArgs(0, 1),
	Run:   ChUser,
}

var userLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all accounts registered",
	Run:   ListUser,
}

var userCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean account db",
	Long:  "Remove all users from inner dbs.",
	Run:   CleanUser,
}

var userRmCmd = &cobra.Command{
	Use:   "remove <UID>",
	Short: "Remove user UID from inner db",
	Args:  cobra.ExactArgs(1),
	Run:   RmUser,
}

var userLookupCmd = &cobra.Command{
	Use:   "lookup <UserName>",
	Short: "Lookup user info by user name",
	Args:  cobra.ExactArgs(1),
	Run:   LookUp,
}

func init() {
	userCmd.AddCommand(userChCmd, userLsCmd, userCleanCmd, userRmCmd, userLookupCmd)
	RootCmd.AddCommand(userCmd)
}

func ChUser(cmd *cobra.Command, params []string) {
	var err error
	var userName string
	if len(params) == 0 {
		userName = ""
	} else {
		userName = params[0]
	}
	err = qshell.ChUser(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chuser: %v\n", err)
		os.Exit(1)
	}
}

func ListUser(cmd *cobra.Command, params []string) {
	err := qshell.ListUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lsuser: %v\n", err)
		os.Exit(1)
	}
}

func CleanUser(cmd *cobra.Command, params []string) {
	err := qshell.CleanUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "CleanUser: %v\n", err)
		os.Exit(1)
	}
}

func RmUser(cmd *cobra.Command, params []string) {
	userName := params[0]
	err := qshell.RmUser(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)
		os.Exit(1)
	}
}

func LookUp(cmd *cobra.Command, params []string) {
	userName := params[0]
	err := qshell.LookUp(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)
		os.Exit(1)
	}
}
