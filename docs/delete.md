# 简介
`delete` 命令用来从七牛的空间里面删除一个文件。

参考文档：[资源删除 (delete)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/delete.html)

# 格式
```
qshell delete <Bucket> <Key>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell delete -h 

// 详细文档（此文档）
$ qshell delete --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间【必选】
- Key：空间中的文件名【必选】             

# 示例
删除空间 `if-pbl` 里面的视频 `qiniu.mp4`
```
qshell delete if-pbl qiniu.mp4
```
