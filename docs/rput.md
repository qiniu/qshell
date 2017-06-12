# 简介

`rput`命令使用七牛支持的分片上传的方式来上传一个文件，一般文件大小较大的情况下，可以使用分片上传来有效地保证文件上传的成功。

参考文档：

[创建块 (mkblk)](http://developer.qiniu.com/code/v6/api/kodo-api/up/mkblk.html)

[上传片 (bput)](http://developer.qiniu.com/code/v6/api/kodo-api/up/bput.html)

[创建文件 (mkfile)](http://developer.qiniu.com/code/v6/api/kodo-api/up/mkfile.html)

# 格式

```
qshell rput <Bucket> <Key> <LocalFile> [Overwrite] [MimeType] [UpHost] [FileType]
```

其中 `Overwrite`，`MimeType`，`UpHost` `FileType`参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

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

1.上传本地文件`/Users/jemy/Documents/qiniu.mp4`到空间`if-pbl`里面。

```
$ qshell rput if-pbl qiniu.mp4 /Users/jemy/Documents/qiniu.mp4
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.mp4 => if-pbl : qiniu.mp4 ...
Progress: 100.00%
Put file /Users/jemy/Documents/qiniu.mp4 => if-pbl : qiniu.mp4 success!
Hash: lhsawSRA9-0L8b0s-cXmojaMhGqn
Fsize: 25538648 ( 24.36 MB )
MimeType: video/mp4
Last time: 16.82 s, Average Speed: 1517.9 KB/s
```

2.上传本地文件`/Users/jemy/Documents/qiniu.mp4`到空间`if-pbl`里面，带前缀`2015/01/18/`，并且指定`MimeType`参数为`video/mp4`

```
$ qshell rput if-pbl 2015/01/18/qiniu.mp4 /Users/jemy/Documents/qiniu.mp4 video/mp4
```
输出：
```
Uploading /Users/jemy/Documents/qiniu.mp4 => if-pbl : 2015/01/18/qiniu.mp4 ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.mp4 => if-pbl : 2015/01/18/qiniu.mp4 success!
Hash: lhsawSRA9-0L8b0s-cXmojaMhGqn
Fsize: 25538648 ( 24.36 MB )
MimeType: video/mp4
Last time: 17.55 s, Average Speed: 1454.9 KB/s
```

3.上传本地文件`/Users/jemy/Documents/qiniu.mp4`到空间`if-pbl`里面，可以指定上传入口`https://upload.qbox.me`。

```
$ qshell rput if-pbl 2015/01/18/qiniu.mp4 /Users/jemy/Documents/qiniu.mp4 https://upload.qbox.me
```

输出:
```
Uploading /Users/jemy/Documents/qiniu.mp4 => if-pbl : 2015/01/18/qiniu.mp4 ...
Progress: 100.00%
Put file /Users/jemy/Documents/qiniu.mp4 => if-pbl : 2015/01/18/qiniu.mp4 success!
Hash: lhsawSRA9-0L8b0s-cXmojaMhGqn
Fsize: 25538648 ( 24.36 MB )
MimeType: video/mp4
Last time: 13.38 s, Average Speed: 1908.6 KB/s
```

5. 使用低频存储

```
$ qshell rput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg true 1
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
