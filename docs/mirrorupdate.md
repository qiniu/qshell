# 简介
`mirrorupdate` 指令用来更新七牛空间中的某个文件。配置了镜像存储的空间，在一个文件首次回源源站拉取资源后，就不再回源了。如果源站更新了一个文件，那么这个文件不会自动被同步更新到七牛空间，这个时候需要使用 `mirrorupdate` 去主动拉取一次这个文件的新内容回来覆盖七牛空间中的旧文件。

功能同 `prefetch`

参考文档：[镜像资源更新 (prefetch)](http://developer.qiniu.com/docs/v6/api/reference/rs/prefetch.html)

# 格式
```
qshell mirrorupdate <Bucket> <Key>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell mirrorupdate -h 

// 详细文档（此文档）
$ qshell mirrorupdate --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或者私有空间【必选】
- Key：空间中文件的名称【必选】

# 示例
```
$ qshell mirrorupdate if-pbl qiniu.png
```