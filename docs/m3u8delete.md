# 简介
`m3u8delete` 命令用来根据七牛空间中 m3u8 播放列表文件名字来删除空间中的 m3u8 播放列表文件和所引用的所有切片文件。

# 格式
```
qshell m3u8delete <Bucket> <M3u8Key>
``` 

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

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
 