# 简介
`forbidden` 禁用指定文件，禁用之后该文件不可访问。

# 格式
```
qshell forbidden <Bucket> <Key> [flags]
```

# 参数
* Bucket：空间名，可以为公开空间或私有空间。【必须】
* Key：空间中的文件名。【必须】

# 选项
* -r/--reverse : 启用指定文件 【可选】

# 示例
```
$ qshell forbidden test-bucket test-file
```
