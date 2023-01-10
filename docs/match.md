# 简介
`match` 指令用来验证本地文件和七牛云存储文件是否匹配。

# 格式
```
qshell match <Bucket> <Key> <LocalFile>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell match -h 

// 详细文档（此文档）
$ qshell match --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或者私有空间【必选】
- Key：空间中文件的名称【必选】
- LocalFile：本地文件路径 【必选】

# 示例
```
$ qshell match if-pbl qiniu.png /Users/lala/Desktop/qiniu.png
```