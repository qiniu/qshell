package cmd

import (
	"fmt"
	"os"

	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var lsCmdEg = `qshell user ls`
var userLsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all user registered",
	Example: lsCmdEg,
	Run:     ListUser,
}

var currentCmdEg = `qshell user current`
var userCurrentCmd = &cobra.Command{
	Use:     "current",
	Short:   "get current user info",
	Example: currentCmdEg,
	Run:     CurrentUser,
}

var lookupCmdEg = `qshell user lookup <UserName>`
var userLookupCmd = &cobra.Command{
	Use:     "lookup <UserName>",
	Short:   "Lookup user info by user name",
	Example: lookupCmdEg,
	Run:     LookUp,
}

var addCmdEg = ` qshell user add <AK> <SK> <UserName>
 or
 qshell user add --ak <AK> --sk <SK> --name <UserName>`
var userAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "add user info to local",
	Example: addCmdEg,
	Run:     Add,
}

var changeUserCmdEg = `qshell user cu <UserName>`
var userChCmd = &cobra.Command{
	Use:     "cu [<UserName>]",
	Short:   "Change user to UserName",
	Example: changeUserCmdEg,
	Run:     ChUser,
}

var cleanCmdEg = `qshell user clean`
var userCleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "clean account db",
	Long:    "Remove all users from inner dbs.",
	Example: cleanCmdEg,
	Run:     CleanUser,
}

var rmCmdEg = `qshell user remove <UserName>`
var userRmCmd = &cobra.Command{
	Use:     "remove <UserName>",
	Short:   "Remove user info from inner db",
	Example: rmCmdEg,
	Run:     RmUser,
}

var (
	userLsName bool
	userAK     string
	userSK     string
	userName   string
)

func init() {
	userLsCmd.Flags().BoolVarP(&userLsName, "name", "n", false, "only list user names")
	userAddCmd.Flags().StringVarP(&userAK, "ak", "", "", "user's access key of Qiniu")
	userAddCmd.Flags().StringVarP(&userSK, "sk", "", "", "user's secret key of Qiniu")
	userAddCmd.Flags().StringVarP(&userName, "name", "", "", "user id of local")
	userCmd.AddCommand(userAddCmd, userRmCmd, userCleanCmd, userChCmd, userLsCmd, userLookupCmd, userCurrentCmd)
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
	err = account.ChUser(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chuser: %v\n", err)
		os.Exit(1)
	}
}

// 列举本地数据库记录的账户
// qshell user ls
func ListUser(cmd *cobra.Command, params []string) {
	err := account.ListUser(userLsName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lsuser: %v\n", err)
		os.Exit(1)
	}
}

// 获取当前账户信息
// qshell user current
func CurrentUser(cmd *cobra.Command, params []string) {
	account, gErr := account.GetAccount()
	if gErr != nil {
		fmt.Println(gErr)
		os.Exit(data.STATUS_ERROR)
	}
	fmt.Println(account.String())
}

// 删除本地记录的数据库
// qshell user clean
func CleanUser(cmd *cobra.Command, params []string) {
	err := account.CleanUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "CleanUser: %v\n", err)
		os.Exit(1)
	}
}

// 删除用户
// qshell user remove <UserName>
func RmUser(cmd *cobra.Command, params []string) {
	if len(params) == 0 {
		log.Error(alert.CannotEmpty("user name", rmCmdEg))
		return
	}
	userName := params[0]
	err := account.RmUser(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RmUser: %v\n", err)
		os.Exit(1)
	}
}

// 查询用用户是否存在本地数据库中
// qshell user lookup <UserName>
func LookUp(cmd *cobra.Command, params []string) {
	if len(params) == 0 {
		log.Error(alert.CannotEmpty("user name", lookupCmdEg))
		return
	}
	userName := params[0]
	err := account.LookUp(userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LookUp: %v\n", err)
		os.Exit(1)
	}
}

// 查询用用户是否存在本地数据库中
// qshell user add <AK> <SK> <UserName>
func Add(cmd *cobra.Command, params []string) {
	if len(params) == 3 {
		userAK = params[0]
		userSK = params[1]
		userName = params[2]
	}

	// check
	if len(userName) == 0 {
		log.Error(alert.CannotEmpty("user name", addCmdEg))
		return
	}
	if len(userAK) == 0 {
		log.Error(alert.CannotEmpty("user ak", addCmdEg))
		return
	}
	if len(userSK) == 0 {
		log.Error(alert.CannotEmpty("user sk", addCmdEg))
		return
	}
	sErr := account.SaveAccount(account.Account{
		Name:      userAK,
		AccessKey: userSK,
		SecretKey: userName,
	}, accountOver)

	if sErr != nil {
		fmt.Println(sErr)
		os.Exit(data.STATUS_ERROR)
	}

}
