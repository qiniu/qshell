# 简介

`fetch`指令根据七牛的公开API [fetch](http://developer.qiniu.com/code/v6/api/kodo-api/rs/fetch.html) 来从互联网上抓取一个资源并保存到七牛的空间中。 
每次抓取的资源，如果指定的Key都是一样的，那么会默认覆盖这个Key所对应的文件。

参考文档：[第三方资源抓取 (fetch)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/fetch.html)

# 格式

```
qshell fetch <RemoteResourceUrl> <Bucket> [<Key>]
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|可选参数|
|-----|-----|------|
|RemoteResourceUrl|互联网上资源的链接，必须是可访问的链接|N|
|Bucket|空间名称，可以为公开空间或者私有空间|N|
|Key|该资源保存在空间中的名字，如果不指定这个名字，那么会使用抓取的资源的内容hash值来作为文件名|Y|

# 示例

1.抓取一个资源并以指定的文件名保存在七牛的空间里面

```
$ qshell fetch https://www.baidu.com/img/bdlogo.png if-pbl bdlogo.png

Key: bdlogo.png
Hash: FrUHIqhkDDd77-AtiDcOwi94YIeM
Fsize: 5331 (5.21 KB)
Mime: image/png

```

2.抓取一个资源并使用文件的内容hash值来作为文件名保存在七牛的空间中

```
$ qshell fetch https://www.baidu.com/img/bdlogo.png if-pbl

Key: FrUHIqhkDDd77-AtiDcOwi94YIeM
Hash: FrUHIqhkDDd77-AtiDcOwi94YIeM
Fsize: 5331 (5.21 KB)
Mime: image/png
```