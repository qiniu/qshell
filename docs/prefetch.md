# 简介

`prefetch`指令用来更新七牛空间中的某个文件。配置了镜像存储的空间，在一个文件首次回源源站拉取资源后，就不再回源了。如果源站更新了一个文件，那么这个文件不会自动被同步更新到七牛空间，这个时候需要使用`prefetch`去主动拉取一次这个文件的新内容回来覆盖七牛空间中的旧文件。

参考文档：[镜像资源更新 (prefetch)](http://developer.qiniu.com/docs/v6/api/reference/rs/prefetch.html) 

# 格式

```
qshell prefetch <Bucket> <Key>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|可选参数|
|-----|-----|------|
|Bucket|空间名称，可以为公开空间或者私有空间|N|
|Key|空间中文件的名称|N|

# 示例

```
$ qshell prefetch if-pbl qiniu.png
```