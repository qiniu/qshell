# 简介
`domains` 指令可以根据指定的空间参数获取和该空间关联的所有域名并输出。

# 格式
```
qshell domains <Bucket> [--detail]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell domains -h 

// 详细文档（此文档）
$ qshell domains --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或者私有空间【必选】

# 选项
--detail：展示域名的详细信息，默认只展示域名的名称。【可选】

# 示例
获取空间 `if-pbl` 对应的所有域名：
```
$ qshell domains if-pbl
```

输出：
```
if-pbl.qiniudn.com
7pn64c.com1.z0.glb.clouddn.com
7pn64c.com2.z0.glb.clouddn.com
7pn64c.com2.z0.glb.qiniucdn.com
```