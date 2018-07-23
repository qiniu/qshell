package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
	"os"
	"strconv"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var userChCmd = &cobra.Command{
	Use:   "cu <AccountName>",
	Short: "Change user to AccountName",
	Args:  cobra.ExactArgs(1),
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
	Run:   ListUser,
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
	uid, err := strconv.Atoi(params[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot convert %s to integer\n", params[0])
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Use)
		os.Exit(1)
	}
	err = qshell.ChUser(uid)
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
	uid, err := strconv.Atoi(params[0])
	if err != nil {
		fmt.Println("user id must be integer")
		os.Exit(1)
	}
	qshell.RmUser(uid)
}

func LookUp(cmd *cobra.Command, params []string) {
	err := qshell.LookUp(params[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)
		os.Exit(1)
	}
}
