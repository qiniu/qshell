# 简介

`delete` 命令用来从七牛的空间里面删除一个文件。

参考文档：[资源删除 (delete)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/delete.html)

# 格式

```
qshell delete <Bucket> <Key>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|	描述|
|------|------|
|Bucket|空间名称，可以为公开空间或者私有空间|
|Key|空间中的文件名|

# 示例

删除空间`if-pbl`里面的视频`qiniu.mp4`

```
qshell delete if-pbl qiniu.mp4
```