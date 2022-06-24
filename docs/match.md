# 简介
`match` 指令用来验证本地文件和七牛云存储文件是否匹配。

# 格式
```
qshell match <Bucket> <Key> <LocalFile>
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数
- Bucket：空间名称，可以为公开空间或者私有空间【必选】
- Key：空间中文件的名称【必选】
- LocalFile：本地文件路径 【必选】

# 示例
```
$ qshell match if-pbl qiniu.png /Users/lala/Desktop/qiniu.png
```