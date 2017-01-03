# 简介

`chgm` 指令用来为空间中的一个文件修改MimeType。

参考文档：[资源元信息修改 (chgm)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/chgm.html)

# 格式
```
qshell chgm <Bucket> <Key> <NewMimeType>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|-----|-----|
|Bucket|空间名称，可以为公开空间或者私有空间|
|Key|空间中的文件名|
|NewMimeType|给文件指定的新的MimeType|

# 示例

修改`if-pbl`空间中`qiniu.png`图片的MimeType为`image/jpeg`

```
$ qshell chgm if-pbl qiniu.png image/jpeg
```

修改完成，我们检查一下文件的MimeType：

```
$ qshell stat if-pbl qiniu.png
```

输出

```
Bucket:             if-pbl
Key:                qiniu.png
Hash:               FrUHIqhkDDd77-AtiDcOwi94YIeM
Fsize:              5331
PutTime:            14285516077733591
MimeType:           image/jpeg

```
我们发现，文件的MimeType已经被修改为`image/jpeg`。