# 简介

`fput`命令用来以`multipart/form-data`的表单方式上传一个文件。适合于中小型文件的上传，一般建议如果文件大小超过100MB的话，都使用分片上传。

参考文档：[直传文件 (upload)](http://developer.qiniu.com/code/v6/api/kodo-api/up/upload.html)

# 格式

```
qshell fput <Bucket> <Key> <LocalFile> [Overwrite] [MimeType] [UpHost] [FileType]
```

其中 `Overwrite`，`MimeType`，`UpHost`，`FileType` 参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名称|描述|可选参数|
|---------|-----------------|----------|
|Bucket|七牛空间名称，可以为公开空间或私有空间|N|
|Key|文件保存在七牛空间的名称|N|
|LocalFile|本地文件的路径|N|
|Overwrite|是否覆盖空间已有文件，默认为`false`|Y|
|MimeType|指定文件的MimeType|Y|
|UpHost|上传入口地址，默认为空间所在机房的上传加速域名|Y|
|FileType|文件存储类型，默认为`0`(标准存储） `1`为低频存储|Y|

关于 `UpHost` ，这个是用来指定上传所使用的入口域名。在不指定的情况下，程序会自动根据空间来获取其所在的机房，并选择对应的上传加速域名作为上传域名。对于七牛的几大机房，默认的上传加速域名和其他源站域名分别如下表。

|机房|上传加速域名|源站上传域名|https上传加速域名|https上传源站域名|
|----|----------------------|--------------------|----------------------|-----------------------|
|华东|http://upload.qiniu.com|http://up.qiniu.com|https://upload.qbox.me|https://up.qbox.me|
|华北|http://upload-z1.qiniu.com|http://up-z1.qiniu.com|https://upload-z1.qbox.me|https://up-z1.qbox.me|
|华南|http://upload-z2.qiniu.com|http://up-z2.qiniu.com|https://upload-z2.qbox.me|https://up-z2.qbox.me|
|北美|http://upload-na0.qiniu.com|http://up-na0.qiniu.com|https://upload-na0.qbox.me|https://up-na0.qbox.me|

当自行指定上传`UpHost`的时候，请根据空间所在机房，从上面的表中选择正确的入口域名。

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
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg image/jpg
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

3.上传本地文件`/Users/jemy/Documents/qiniu.jpg`到空间`if-pbl`里面，并且指定指定的上传入口

```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg https://upload.qbox.me
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpeg
Last time: 1.47 s, Average Speed: 26.7 KB/s
```

4.覆盖上传 `qiniu.mp4' 到空间`if-pbl`

```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg true
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
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg true 1
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
