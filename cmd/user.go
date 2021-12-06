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

// 列举本地数据库记录的账户
var lsCmdEg = `qshell user ls`
var userLsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all user registered",
	Example: lsCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		err := account.ListUser(userLsName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lsuser: %v\n", err)
			os.Exit(1)
		}
	},
}

// 获取当前账户信息
var currentCmdEg = `qshell user current`
var userCurrentCmd = &cobra.Command{
	Use:     "current",
	Short:   "get current user info",
	Example: currentCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		account, gErr := account.GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(data.STATUS_ERROR)
		}
		fmt.Println(account.String())
	},
}

// 查询用用户是否存在本地数据库中
var lookupCmdEg = `qshell user lookup <UserName>`
var userLookupCmd = &cobra.Command{
	Use:     "lookup <UserName>",
	Short:   "Lookup user info by user name",
	Example: lookupCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error(alert.CannotEmpty("user name", lookupCmdEg))
			return
		}
		userName := args[0]
		err := account.LookUp(userName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LookUp: %v\n", err)
			os.Exit(1)
		}
	},
}

// 添加用户
var addCmdEg = ` qshell user add <AK> <SK> <UserName>
 or
 qshell user add --ak <AK> --sk <SK> --name <UserName>`
var userAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "add user info to local",
	Example: addCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 3 {
			userAK = args[0]
			userSK = args[1]
			userName = args[2]
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
		}, userOverwrite)

		if sErr != nil {
			fmt.Println(sErr)
			os.Exit(data.STATUS_ERROR)
		}
	},
}

// 切换用户
var changeUserCmdEg = `qshell user cu <UserName>`
var userChCmd = &cobra.Command{
	Use:     "cu [<UserName>]",
	Short:   "Change user to UserName",
	Example: changeUserCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var userName string
		if len(args) == 0 {
			userName = ""
		} else {
			userName = args[0]
		}
		err = account.ChUser(userName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "chuser: %v\n", err)
			os.Exit(1)
		}
	},
}

// 删除本地记录的数据库
var cleanCmdEg = `qshell user clean`
var userCleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "clean account db",
	Long:    "Remove all users from inner dbs.",
	Example: cleanCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		err := account.CleanUser()
		if err != nil {
			fmt.Fprintf(os.Stderr, "CleanUser: %v\n", err)
			os.Exit(1)
		}
	},
}

// 删除用户
var rmCmdEg = `qshell user remove <UserName>`
var userRmCmd = &cobra.Command{
	Use:     "remove <UserName>",
	Short:   "Remove user info from inner db",
	Example: rmCmdEg,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error(alert.CannotEmpty("user name", rmCmdEg))
			return
		}
		userName := args[0]
		err := account.RmUser(userName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "RmUser: %v\n", err)
			os.Exit(1)
		}
	},
}

var (
	userLsName bool

	userAK        string
	userSK        string
	userName      string
	userOverwrite bool
)

func init() {
	userLsCmd.Flags().BoolVarP(&userLsName, "name", "n", false, "only list user names")

	userAddCmd.Flags().StringVarP(&userAK, "ak", "", "", "user's access key of Qiniu")
	userAddCmd.Flags().StringVarP(&userSK, "sk", "", "", "user's secret key of Qiniu")
	userAddCmd.Flags().StringVarP(&userName, "name", "", "", "user id of local")
	userAddCmd.Flags().BoolVarP(&userOverwrite, "overwrite", "w", false, "overwrite user or not when account exists in local db, by default not overwrite")

	userCmd.AddCommand(
		userAddCmd,     // 添加用户
		userRmCmd,      // 删除某个用户
		userCleanCmd,   // 清除所有用户
		userChCmd,      // 切换当前账户
		userLsCmd,      // 列举所有用户信息
		userLookupCmd,  // 查看某个用户的信息
		userCurrentCmd, // 查看当前用户信息
	)
	RootCmd.AddCommand(userCmd)
}
