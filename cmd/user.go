package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/account/operations"
	"github.com/spf13/cobra"
)

var accountCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {

	var accountOver bool
	var cmd = &cobra.Command{
		Use:   "account [<AccessKey> <SecretKey> <Name>]",
		Short: "Get/Set current account's Id and SecretKey",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.Account
			if len(args) == 0 {
				operations.Current(cfg)
				return
			}

			info := operations.AddInfo{
				Over: accountOver,
			}
			if len(args) > 0 {
				info.AccessKey = args[0]
			}
			if len(args) > 1 {
				info.SecretKey = args[1]
			}
			if len(args) > 2 {
				info.Name = args[2]
			}
			operations.Add(cfg, info)
		},
	}

	cmd.Flags().BoolVarP(&accountOver, "overwrite", "w", false, "overwrite account or not when account exists in local db, by default not overwrite")

	return cmd
}

var userCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			operations.User(cfg, operations.UserInfo{})
		},
	}
	return cmd
}

// 列举本地数据库记录的账户
var userLsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ListInfo{}
	var cmd = &cobra.Command{
		Use:     "ls",
		Short:   "List all user registered",
		Example: `qshell user ls`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			operations.List(cfg, info)
		},
	}

	cmd.Flags().BoolVarP(&info.OnlyListName, "name", "n", false, "only list user names")

	return cmd
}

// 获取当前账户信息
var userCurrentCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "current",
		Short:   "get current user info",
		Example: `qshell user current`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			operations.Current(cfg)
		},
	}
	return cmd
}

// 查询用用户是否存在本地数据库中
var userLookupCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.LookUpInfo{}
	var cmd = &cobra.Command{
		Use:     "lookup <UserName>",
		Short:   "Lookup user info by user name",
		Example: `qshell user lookup <UserName>`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			if len(args) > 0 {
				info.Name = args[0]
			}
			operations.LookUp(cfg, info)
		},
	}
	return cmd
}

// 添加用户
var userAddCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.AddInfo{}
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "add user info to local",
		Example: `qshell user add <AK> <SK> <UserName>
 or
 qshell user add --ak <AK> --sk <SK> --name <UserName>`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			if len(args) == 3 {
				info.AccessKey = args[0]
				info.SecretKey = args[1]
				info.Name = args[2]
			}
			operations.Add(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.AccessKey, "ak", "", "", "user's access key of Qiniu")
	cmd.Flags().StringVarP(&info.SecretKey, "sk", "", "", "user's secret key of Qiniu")
	cmd.Flags().StringVarP(&info.Name, "name", "", "", "user id of local")
	cmd.Flags().BoolVarP(&info.Over, "overwrite", "w", false, "overwrite user or not when account exists in local db, by default not overwrite")

	return cmd
}

// 切换用户
var userChCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.ChangeInfo{}
	var cmd = &cobra.Command{
		Use:     "cu [<UserName>]",
		Short:   "Change user to UserName",
		Example: `qshell user cu <UserName>`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			if len(args) > 0 {
				info.Name = args[0]
			}
			operations.Change(cfg, info)
		},
	}
	return cmd
}

// 删除本地记录的数据库
var userCleanCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "clean",
		Short:   "clean account db",
		Long:    "Remove all users from inner dbs.",
		Example: `qshell user clean`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			operations.Clean(cfg, operations.CleanInfo{})
		},
	}
	return cmd
}

// 删除用户
var userRmCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var info = operations.RemoveInfo{}
	var cmd = &cobra.Command{
		Use:     "remove <UserName>",
		Short:   "Remove user info from inner db",
		Example: `qshell user remove <UserName>`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.User
			if len(args) > 0 {
				info.Name = args[0]
			}
			operations.Remove(cfg, info)
		},
	}
	return cmd
}

func init() {
	registerLoader(userCmdLoader)
}

func userCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	userCmd := userCmdBuilder(cfg)
	userCmd.AddCommand(
		userAddCmdBuilder(cfg),     // 添加用户
		userRmCmdBuilder(cfg),      // 删除某个用户
		userCleanCmdBuilder(cfg),   // 清除所有用户
		userChCmdBuilder(cfg),      // 切换当前账户
		userLsCmdBuilder(cfg),      // 列举所有用户信息
		userLookupCmdBuilder(cfg),  // 查看某个用户的信息
		userCurrentCmdBuilder(cfg), // 查看当前用户信息
	)

	superCmd.AddCommand(
		accountCmdBuilder(cfg),
		userCmd,
	)
}
