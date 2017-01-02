# 简介

`stat`指令根据七牛的公开API [stat](http://developer.qiniu.com/code/v6/api/kodo-api/rs/stat.html) 来获取空间中的一个文件的基本信息，包括文件的名称，保存的时间，hash值，文件的大小和MimeType。

参考文档：[资源元信息查询 (stat)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/stat.html)

# 格式

```
qshell stat <Bucket> <Key>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|-----|-----|
|Bucket|空间名称，可以为公开空间或者私有空间|
|Key|空间中的文件名|


# 示例

获取空间`if-pbl`中文件`qiniu.png`的基本信息

```
$ qshell stat if-pbl qiniu.png
```

输出：

```
Bucket:             if-pbl
Key:                test.png
Hash:               FpHyF0kkil3sp-SaXXX8TBJY3jDh
Fsize:              163469 -> 159.64 KB
PutTime:            14832477937031429 -> 2017-01-01 13:16:33.7031429 +0800 CST
MimeType:           image/png
```