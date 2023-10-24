# 简介
`expire` 指令用来为空间中的一个文件修改 **过期时间**。（即多长时间后，该文件会被自动删除）

## 注：
v2.7.0 到 v2.9.2 不支持取消过期时间配置，如果过期天数配置为 0 会立即删除文件。
v2.10.0 开始，过期天数配置为 0 时为取消过期时间配置。

# 格式
```
qshell expire <Bucket> <Key> <DeleteAfterDays>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell expire -h 

// 详细文档（此文档）
$ qshell expire --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或者私有空间【必选】
- Key：空间中的文件名【必选】
- DeleteAfterDays：给文件指定的新过期时间，范围：大于等于 0，0：取消过期时间设置。单位为：天【必选】

# 示例
修改 `if-pbl` 空间中 `qiniu.png` 图片的过期时间为：`3天后自动删除`
```
$ qshell expire if-pbl qiniu.png 3
```
输入该命令，后该文件就已经被修改为 3 天后自动删除了。
