# 简介
`m3u8delete` 命令用来根据七牛空间中 m3u8 播放列表文件名字来删除空间中的 m3u8 播放列表文件和所引用的所有切片文件。

# 格式
```
qshell m3u8delete <Bucket> <M3u8Key>
``` 

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell m3u8delete -h 

// 详细文档（此文档）
$ qshell m3u8delete --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：m3u8 文件所在空间，可以为公开空间或私有空间 【必选】
- M3u8Key：m3u8 文件的名字 【必选】

# 示例
1 删除公开空间中 m3u8 文件及其所引用的所有切片文件。
```
qshell m3u8delete if-pbl qiniu.m3u8
```

2 删除私有空间中 m3u8 文件及其所引用的所有切片文件。
```
qshell m3u8delete if-pri qiniu.m3u8
```
 