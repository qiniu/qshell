# 简介
`user` 命令用来对本地数据库中存储的账户信息进行管理，可以添加账号、查看/切换当前账号、列举本地保存的账号、移除特定的账号。

# 格式
```
qshell user <子命令>
``` 

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell user -h 

// 详细文档（此文档）
$ qshell user --doc

// 子命令示例
$ qshell user cu -h
$ qshell user cu --doc
```

# 鉴权
无


# 子命令
user的字命令有：
* add：添加账号
* clean：清除本地数据库
* cu：切换当前的账户
* current：查看当前账号
* lookup：通过用户名字查找用户信息
* ls：列出所有本地的账户信息
* remove：移除特定用户

# 示例
1. 添加账号
```
//  --ak：七牛账号对应的 AccessKey [获取](https://portal.qiniu.com/user/key) 。【必选】
//  --sk：七牛账号对应的 SecretKey [获取](https://portal.qiniu.com/user/key) 。【必选】
//  --name：AccessKey 和 SecretKey 对的 id, 可以任意取，但同一台机器此 id 不可重复；和在七牛注册的邮箱信息没有关系， 只是 qshell 本地用来标示 <ak, sk> 对。【必选】

qshell user add --ak ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6x --sk LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKiw --name name_test
```

2. 清除本地数据库
``` 
qshell user clean // 注：仅仅清除本地数据库，会保留当前账户
```

3. 切换当前的账户
```
qshell user cu       // 切换至上一次使用的账户
qshell user cu test  // 切换到 `test` 账号，`test` 为 ak,sk 对的 id
```

4. 列举本地所有的账号信息
```
qshell user ls
```

5. 输出某个账户信息
```
qshell user lookup test // `test` 为 ak,sk 对的 id
```

6. 删除 `test` 账号
```
qshell user remove test // `test` 为 ak,sk 对的 id
```
