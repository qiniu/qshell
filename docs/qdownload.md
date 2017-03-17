# 简介

`qdownload`可以将七牛空间中的文件同步到本地磁盘中。支持只同步带特定前缀或者后缀的文件，也支持在本地备份路径不变的情况下进行增量同步。
需要额外指出的是，将文件同步到本地都是走的七牛存储源站的流量而不是CDN的流量，因为CDN通常情况下会认为大量的文件下载操作是非法访问从而进行限制。

本工具批量下载文件支持多文件并发下载，另外还支持单个文件的断点续传。除此之外，也可以支持指定前缀或者后缀的文件同步，注意这里的前缀只能指定一个，但是后缀可以指定多个，多个后缀直接使用英文的逗号(,)分隔。

# 格式

```
qshell qdownload [<ThreadCount>] <LocalDownloadConfig>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名称|描述|可选参数|取值范围|
|----------|-----------|----------|---------|
|ThreadCount|下载的并发协程数量|Y|1-2000，如果不在这个范围内，默认为5|
|LocalDownloadConfig|本地下载的配置文件，内容包括要下载的文件所在空间，文件前缀等信息，具体参考配置文件说明|N||

其中 `ThreadCount` 表示支持同时下载多个文件。

# 配置

`qdownload` 功能需要配置文件的支持，配置文件的内容如下：

```
{
    "dest_dir"   :   "<LocalBackupDir>",
    "bucket"     :   "<Bucket>",
    "prefix"     :   "image/",
    "suffixes"   :   ".png,.jpg"
}
```

|参数名|描述|可选参数|
|--------------|---------------|----------------|
|dest_dir|本地数据备份路径，为全路径|N|
|bucket|空间名称|N|
|prefix|只同步指定前缀的文件，默认为空|Y|
|suffixes|只同步指定后缀的文件，默认为空|Y|
|cdn_domain|指定下载的域名，如果不设置则默认使用源站下载|Y|


注意，在Windows系统下面使用的时候，注意`dest_dir`的设置遵循`D:\\jemy\\backup`这种方式。也就是路径里面的`\`要有两个（`\\`）。

# 示例

需要同步空间`qdisk`中的所有以`movies/`开头，并以`.mp4`结尾的文件到本地路径`/Users/jemy/Temp7/backup`下面：

配置文件：

```
{
	"dest_dir"	:	"/Users/jemy/Temp7/backup",
	"bucket"	:	"qdisk",
	"prefix"	:	"movies/",
	"suffixes"	:	".mp4"
}
```

运行命令（下载并发数表示可以同时下载10个文件）：

```
qshell qdownload 10 qdisk_down.conf
```
