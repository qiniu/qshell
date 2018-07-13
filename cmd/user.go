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
	Use:   "cu",
	Short: "Change user to AccountName",
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
	Use:   "remove",
	Short: "Remove users from inner db",
	Run:   RmUser,
}

var userLookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup user info by uid",
	Run:   LookUp,
}

func init() {
	userCmd.AddCommand(userChCmd, userLsCmd, userCleanCmd, userRmCmd, userLookupCmd)
}

func ChUser(cmd *cobra.Command, params []string) {
	if len(params) > 1 || len(params) <= 0 {
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Use)
		os.Exit(1)
	}
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
	if len(params) != 0 {
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Use)
		os.Exit(1)
	}
	err := qshell.ListUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lsuser: %v\n", err)
		os.Exit(1)
	}
}

func CleanUser(cmd *cobra.Command, params []string) {
	if len(params) != 0 {
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Use)
		os.Exit(1)
	}
	err := qshell.CleanUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "CleanUser: %v\n", err)
		os.Exit(1)
	}
}

func RmUser(cmd *cobra.Command, params []string) {
	uids := make([]int, len(params))
	if len(params) > 0 {
		for ind, u := range params {
			uid, err := strconv.Atoi(u)
			if err != nil {
				fmt.Println("user ids must be integer")
				os.Exit(1)
			}
			uids[ind] = uid
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Use)
	}
	qshell.RmUser(uids)
}

func LookUp(cmd *cobra.Command, params []string) {
	if len(params) != 1 {
		fmt.Fprintf(os.Stderr, "lookup need a username\n")
		fmt.Fprintf(os.Stderr, "%s\n", cmd.Use)
		os.Exit(1)
	}
	err := qshell.LookUp(params[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", cmd, err)
		os.Exit(1)
	}
}
