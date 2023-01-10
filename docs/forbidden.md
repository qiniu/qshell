# 简介
`forbidden` 禁用指定文件，禁用之后该文件不可访问。

# 格式
```
qshell forbidden <Bucket> <Key> [flags]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell forbidden -h 

// 详细文档（此文档）
$ qshell forbidden --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
* Bucket：空间名，可以为公开空间或私有空间。【必须】
* Key：空间中的文件名。【必须】

# 选项
* -r/--reverse : 启用指定文件 【可选】

# 示例
```
$ qshell forbidden test-bucket test-file
```
