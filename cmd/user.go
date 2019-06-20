package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
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
	Use:   "remove <UserName>",
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

var (
	userLsName bool
)

func init() {
	userLsCmd.Flags().BoolVarP(&userLsName, "name", "n", false, "only list user names")
	userCmd.AddCommand(userChCmd, userLsCmd, userCleanCmd, userRmCmd, userLookupCmd)
	RootCmd.AddCommand(userCmd)
}

// 切换用户
// qshell user cu <Name>
func ChUser(cmd *cobra.Command, params []string) {
	var err error
	var userName string
	if len(params) == 0 {
		userName = ""
	} else {
		userName = params[0]
	}
	err = iqshell.ChUser(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chuser: %v\n", err)
		os.Exit(1)
	}
}

// 列举本地数据库记录的账户
// qshell user ls
func ListUser(cmd *cobra.Command, params []string) {
	err := iqshell.ListUser(userLsName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lsuser: %v\n", err)
		os.Exit(1)
	}
}

// 删除本地记录的数据库
// qshell user clean
func CleanUser(cmd *cobra.Command, params []string) {
	err := iqshell.CleanUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "CleanUser: %v\n", err)
		os.Exit(1)
	}
}

// 删除用户
// qshell user remove <UserName>
func RmUser(cmd *cobra.Command, params []string) {
	userName := params[0]
	err := iqshell.RmUser(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RmUser: %v\n", err)
		os.Exit(1)
	}
}

// 查询用用户是否存在本地数据库中
// qshell user lookup <UserName>
func LookUp(cmd *cobra.Command, params []string) {
	userName := params[0]
	err := iqshell.LookUp(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LookUp: %v\n", err)
		os.Exit(1)
	}
}
