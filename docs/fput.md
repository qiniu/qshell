# 简介

`fput`命令用来以`multipart/form-data`的表单方式上传一个文件。适合于中小型文件的上传，一般建议如果文件大小超过100MB的话，都使用分片上传。

参考文档：[直传文件 (upload)](http://developer.qiniu.com/code/v6/api/kodo-api/up/upload.html)

# 格式

```
qshell fput [--overwrite] [--storage <StorageType>] [--mimetype <MimeType>] <Bucket> <Key> <LocalFile>
```

其中 `Overwrite`，`MimeType`，`StorageType` 参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey`和`Name`的情况下使用。

# 参数

|参数名称|描述|可选参数|
|---------|-----------------|----------|
|Bucket|七牛空间名称，可以为公开空间或私有空间|N|
|Key|文件保存在七牛空间的名称|N|
|LocalFile|本地文件的路径|N|
|Overwrite|是否覆盖空间已有文件，默认为`false`|Y|
|MimeType|指定文件的MimeType|Y|
|StorageType|文件存储类型，默认为`0`(标准存储） `1`为低频存储|Y|

# 示例

1.上传本地文件`/Users/jemy/Documents/qiniu.jpg`到空间`if-pbl`里面。

```
$ qshell fput if-pbl qiniu.jpg /Users/jemy/Documents/qiniu.jpg
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpeg
Last time: 0.33 s, Average Speed: 118.6 KB/s
```

2.上传本地文件`/Users/jemy/Documents/qiniu.jpg`到空间`if-pbl`里面，带前缀`2015/01/18/`，并且指定`MimeType`参数为`image/jpg`

```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --mimetype image/jpg
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpg
Last time: 0.39 s, Average Speed: 101.4 KB/s
```

3.覆盖上传 `qiniu.mp4' 到空间`if-pbl`

```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --overwrite
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpeg
Last time: 0.40 s, Average Speed: 98.2 KB/s
```


5. 使用低频存储

```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --storage 1
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpeg
Last time: 0.40 s, Average Speed: 98.2 KB/s
```
