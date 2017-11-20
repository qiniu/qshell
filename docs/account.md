# 简介

`account`命令用来设置当前用户的`AccessKey`和`SecretKey`，这对Key主要用在其他的需要授权的命令中，比如`stat`,`delete`,`listbucket`命令中。
该命令设置的信息，经过加密保存在命令执行的目录下的`.qshell/account.json`文件中。

# 格式

```
qshell account
``` 

打印当前设置的`AccessKey`和`SecretKey`

```
qshell account <Your AccessKey> <Your SecretKey>
``` 

设置当前用户的`AccessKey`和`SecretKey`

# 参数

|参数名|描述|
|--------|--------|
|AccessKey|七牛账号对应的AccessKey [获取](https://portal.qiniu.com/user/key)|
|SecretKey|七牛账号对应的SecretKey [获取](https://portal.qiniu.com/user/key)|

# 示例

1.设置当前用户的AccessKey和SecretKey

```
qshell account ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6x LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKiw
```

2.输出当前用户设置的AccessKey和SecretKey

```
qshell account
```
输出:

```
AccessKey: ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6x
SecretKey: LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKiw
```
