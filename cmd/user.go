package cmd

import (
	"fmt"
	"os"

	"github.com/qiniu/qshell/v2/iqshell/common/account"
	"github.com/qiniu/qshell/v2/iqshell/common/account/operations"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/spf13/cobra"
)

var userCmdBuilder = func() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}
	return cmd
}

// 列举本地数据库记录的账户
var userLsCmdBuilder = func() *cobra.Command {
	var name = false
	var cmdEg = `qshell user ls`
	var cmd = &cobra.Command{
		Use:     "ls",
		Short:   "List all user registered",
		Example: cmdEg,
		Run: func(cmd *cobra.Command, args []string) {
			err := account.ListUser(name)
			if err != nil {
				log.Error("lsuser: %v", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&name, "name", "n", false, "only list user names")

	return cmd
}

// 获取当前账户信息
var userCurrentCmdBuilder = func() *cobra.Command {
	var cmdEg = `qshell user current`
	var cmd = &cobra.Command{
		Use:     "current",
		Short:   "get current user info",
		Example: cmdEg,
		Run: func(cmd *cobra.Command, args []string) {
			account, gErr := account.GetAccount()
			if gErr != nil {
				log.Error("user current: %v", gErr)
				os.Exit(data.STATUS_ERROR)
			}
			log.Alert(account.String())
		},
	}
	return cmd
}

// 查询用用户是否存在本地数据库中
var userLookupCmdBuilder = func() *cobra.Command {
	var cmdEg = `qshell user lookup <UserName>`
	var cmd = &cobra.Command{
		Use:     "lookup <UserName>",
		Short:   "Lookup user info by user name",
		Example: cmdEg,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Error(alert.CannotEmpty("user name", cmdEg))
				return
			}
			userName := args[0]
			err := account.LookUp(userName)
			if err != nil {
				log.Error("LookUp: %v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

// 添加用户
var userAddCmdBuilder = func() *cobra.Command {

	var cmdInfo = operations.AddInfo{}
	var addCmdEg = ` qshell user add <AK> <SK> <UserName>
 or
 qshell user add --ak <AK> --sk <SK> --name <UserName>`

	var cmd = &cobra.Command{
		Use:     "add",
		Short:   "add user info to local",
		Example: addCmdEg,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 3 {
				cmdInfo.AccessKey = args[0]
				cmdInfo.SecretKey = args[1]
				cmdInfo.Name = args[2]
			}

			// check
			if len(cmdInfo.Name) == 0 {
				log.Error(alert.CannotEmpty("user name", addCmdEg))
				return
			}
			if len(cmdInfo.AccessKey) == 0 {
				log.Error(alert.CannotEmpty("user ak", addCmdEg))
				return
			}
			if len(cmdInfo.SecretKey) == 0 {
				log.Error(alert.CannotEmpty("user sk", addCmdEg))
				return
			}

			operations.Add(cmdInfo)
		},
	}

	cmd.Flags().StringVarP(&cmdInfo.AccessKey, "ak", "", "", "user's access key of Qiniu")
	cmd.Flags().StringVarP(&cmdInfo.SecretKey, "sk", "", "", "user's secret key of Qiniu")
	cmd.Flags().StringVarP(&cmdInfo.Name, "name", "", "", "user id of local")
	cmd.Flags().BoolVarP(&cmdInfo.Over, "overwrite", "w", false, "overwrite user or not when account exists in local db, by default not overwrite")

	return cmd
}

// 切换用户
var userChCmdBuilder = func() *cobra.Command {
	var changeUserCmdEg = `qshell user cu <UserName>`
	var cmd = &cobra.Command{
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
	return cmd
}

// 删除本地记录的数据库
var userCleanCmdBuilder = func() *cobra.Command {
	var cleanCmdEg = `qshell user clean`
	var cmd = &cobra.Command{
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
	return cmd
}

// 删除用户
var userRmCmdBuilder = func() *cobra.Command {
	var rmCmdEg = `qshell user remove <UserName>`
	var cmd = &cobra.Command{
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
	return cmd
}

func init() {
	userCmd := userCmdBuilder()
	userCmd.AddCommand(
		userAddCmdBuilder(),     // 添加用户
		userRmCmdBuilder(),      // 删除某个用户
		userCleanCmdBuilder(),   // 清除所有用户
		userChCmdBuilder(),      // 切换当前账户
		userLsCmdBuilder(),      // 列举所有用户信息
		userLookupCmdBuilder(),  // 查看某个用户的信息
		userCurrentCmdBuilder(), // 查看当前用户信息
	)
	RootCmd.AddCommand(userCmd)
}
