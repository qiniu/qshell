# 简介

`m3u8delete`命令用来根据七牛空间中m3u8播放列表文件名字来删除空间中的m3u8播放列表文件和所引用的所有切片文件。

# 格式

```
qshell m3u8delete <Bucket> <M3u8Key>
``` 

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|可选参数|
|--------|--------|-------|
|Bucket|m3u8文件所在空间，可以为公开空间或私有空间|N|
|M3u8Key|m3u8文件的名字|N|

# 示例

1.删除公开空间中m3u8文件及其所引用的所有切片文件。

```
qshell m3u8delete if-pbl qiniu.m3u8
```

2.删除私有空间中m3u8文件及其所引用的所有切片文件。

```
qshell m3u8delete if-pri qiniu.m3u8
```
 