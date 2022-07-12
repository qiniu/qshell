# 简介
`account` 命令用来设置当前用户的 `AccessKey` 和 `SecretKey` ，这对 Key 主要用在其他的需要授权的命令中，比如 `stat` , `delete` , `listbucket` 命令中。
该命令设置的信息，经过加密保存在 HOME 目录下的 `.qshell/account.json` 文件中。

本地数据库会记录 `account` 注册的所有 <AccessKey> ,  <SecretKey>  和 <Name> 的信息， 所以当用 `account` 注册账户信息时，如果 qshell发 现本地数据库有同样的名字为
<Name> 的账户， 那么默认 qshell 会返回错误信息报告该名字的账户已经存在，如果要覆盖注册，需要使用强制覆盖选项 --overwrite 或者 -w

# 格式
```
qshell account
``` 

打印当前设置的`AccessKey`, `SecretKey` 和 `Name`
```
qshell account [--overwrite | -w]<Your AccessKey> <Your SecretKey> <Your Account Name>
``` 

设置当前用户的`AccessKey`, `SecretKey`和`Name`, Name是用户可以任意取的名字，表示当前在本地记录的账户的名称，和在七牛注册的邮箱信息没有关系。如果 `AccessKey`, `SecretKey`, `Name` 首字母是 "-" , 需要使用` $ qshell account -- ak sk name`的方式添加账号, 这样避免把该项识别成命令行选项。

# 参数
- AccessKey：七牛账号对应的 AccessKey [获取](https://portal.qiniu.com/user/key) 。【必选】
- SecretKey：七牛账号对应的 SecretKey [获取](https://portal.qiniu.com/user/key) 。【必选】
- Name：AccessKey 和 SecretKey 对的 id, 可以任意取，但同一台机器此 id 不可重复；和在七牛注册的邮箱信息没有关系， 只是 qshell 本地用来标示 <ak, sk> 对。【必选】

# 选项
-w/--overwrite: 强制覆盖已经存在的账户

# 示例
1 设置当前用户的 AccessKey, SecretKey, Name
```
qshell account ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6x LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKiw name_test
```

2 输出当前用户设置的 AccessKey 和 SecretKey
```
qshell account
```
输出:
```
Name: name_test
AccessKey: ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6x
SecretKey: LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKiw
```

3 我们可以在设置name_test账户后，继续添加一个账户
```
qshell account ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6abc LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDthaha name_test2
```
qshell 可以记录多个设置的账户信息，账户的管理，切换，删除等，可以参考 qshell user 自命令[文档](user.md)
