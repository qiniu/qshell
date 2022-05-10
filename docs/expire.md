# 简介

`expire` 指令用来为空间中的一个文件修改 **过期时间**。（即多长时间后，该文件会被自动删除）


# 格式
```
qshell expire <Bucket> <Key> <DeleteAfterDays>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|-----|-----|
|Bucket|空间名称，可以为公开空间或者私有空间|
|Key|空间中的文件名|
|DeleteAfterDays |给文件指定的新过期时间，单位为：天|

# 示例

修改`if-pbl`空间中`qiniu.png`图片的过期时间为：`3天后自动删除`

```
$ qshell chtype if-pbl qiniu.png 3
```
输入该命令，后该文件就已经被修改为3天后自动删除了
