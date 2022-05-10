# 简介
`qdownload` 可以将七牛空间中的文件同步到本地磁盘中。支持只同步带特定前缀或者后缀的文件，也支持在本地备份路径不变的情况下进行增量同步。 需要额外指出的是，将文件同步到本地都是走的七牛存储源站的流量而不是 `CDN`
的流量，因为 `CDN` 通常情况下会认为大量的文件下载操作是非法访问从而进行限制。

### 注：【该功能默认需要计费，如果希望享受 10G 的免费流量，请自行设置 cdn_domain 参数，如不设置，需支付源站流量费用，无法减免！！！】

本工具批量下载文件支持多文件并发下载，另外还支持单个文件的断点续传。除此之外，也可以支持指定前缀或者后缀的文件同步，注意这里的前缀只能指定一个，但是后缀可以指定多个，多个后缀直接使用英文的逗号(,)分隔。

# 格式
```
qshell qdownload [-c <ThreadCount>] <LocalDownloadConfig>
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name` 的情况下使用。

# 参数
- LocalDownloadConfig：本地下载的配置文件，内容包括要下载的文件所在空间，文件前缀等信息，具体参考配置文件说明 【必选】

其中 `ThreadCount` 表示支持同时下载多个文件。

# 选项
- -c/--thread：配置下载的并发协程数量，表示支持同时下载多个文件（ThreadCount）, 大小必须在1-2000，如果不在这个范围内，默认为5。

`qdownload` 功能需要配置文件的支持，配置文件的内容如下：
```
{
    "dest_dir"               :   "<LocalBackupDir>",
    "bucket"                 :   "<Bucket>",
    "prefix"                 :   "image/",
    "suffixes"               :   ".png,.jpg",
    "key_file"               :   "<KeyFile>",
    "check_hash"             :   false,
    "cdn_domain"             :   "down.example.com",
    "referer"                :   "http://www.example.com",
    "public"                 :   true,
    "remove_temp_while_error": false,
    "log_file"               :   "download.log",
    "log_level"              :   "info",
    "log_rotate"             :   1,
    "log_stdout"             :   false
}
```

字段说明：

- dest_dir：本地数据备份路径，为全路径 【必选】
- bucket：空间名称 【必选】
- prefix：只同步指定前缀的文件，默认为空 【可选】
- suffixes：只同步指定后缀的文件，默认为空 【可选】
- key_file：配置一个文件，指定需要下载的 keys；默认为空，全量下载 bucket 中的文件 【可选】
- check_hash：是否验证 hash，如果开启可能会耗费较长时间，默认为 `false` 【可选】
- cdn_domain：设置下载的 CDN 域名，默认为空表示从存储源站下载，【该功能默认需要计费，如果希望享受 10G 的免费流量，请自行设置 cdn_domain 参数，如不设置，需支付源站流量费用，无法减免！！！】 【可选】
- referer：如果 CDN 域名配置了域名白名单防盗链，需要指定一个允许访问的 referer 地址；默认为空 【可选】
- public：空间是否为公开空间；为 `true` 时为公有空间，公有空间下载时不会对下载 URL 进行签名，可以提升 CDN 域名性能，默认为 `false`（私有空间）【可选】
- remove_temp_while_error: 当下载遇到错误时删除之前下载的部分文件缓存，默认为 `false` (不删除)【可选】
- log_level：下载日志输出级别，可选值为 `debug`,`info`,`warn`,`error`，其他任何字段均会导致不输出日志。默认 `debug` 。【可选】
- log_file：下载日志的输出文件，默认为输出到 `record_root` 指定的文件中，具体文件路径可以在终端输出看到。【可选】
- log_rotate：下载日志文件的切换周期，单位为天，默认为 7 天即切换到新的下载日志文件 【可选】
- log_stdout：下载日志是否同时输出一份到标准终端，默认为 `false`，主要在调试下载功能时可以指定为 `true` 【可选】
- record_root：下载记录信息保存路径，包括日志文件和下载进度文件；默认为 `qshell` 下载目录；【可选】
    - 通过 `-L` 指定工作目录时，`record_root` 则为 `此工作目录/qupload/$jobId`，
    - 未通过 `-L` 指定工作目录时为 `用户目录/.qshell/users/$CurrentUserName/qupload/$jobId`
    - 注意 `jobId` 是根据上传任务动态生成；据图方式为 MD5("DestDir:$Bucket:KeyFile")；`CurrentUserName` 当前用户的名称

##### 备注：

1. 在Windows系统下面使用的时候，注意 `dest_dir` 的设置遵循 `D:\\jemy\\backup` 这种方式。也就是路径里面的 `\` 要有两个（`\\`）。
2. 在默认不指定 `cdn_domain` 的情况下，会从存储源站下载资源，这部分下载产生的流量会生成存储源站下载流量的计费，请注意，这部分计费不在七牛 CDN 免费 10G 流量覆盖范围。

# 示例

需要同步空间 `qdisk` 中的所有以 `movies/` 开头(理解为前缀的概念，那么 `movies/1.mp4`, `movies/2.mp4` 等以 `movies/` 为前缀的文件都会被下载保存)，并以 `.mp4`
结尾的文件到本地路径 `/Users/jemy/Temp7/backup` 下面（把下面的配置内容写入配置文件 `qdisk_down.conf`，该配置文件需要自行创建）：

```
{
	"dest_dir"	:	"/Users/jemy/Temp7/backup",
	"bucket"	:	"qdisk",
	"cdn_domain"    :      "if-pbl.qiniudn.com",
	"prefix"	:	"movies/",
	"suffixes"	:	".mp4",
	"check_hash"    : false
}
```

运行命令（下载并发数表示可以同时下载 10 个文件）：

```
qshell qdownload -c 10 qdisk_down.conf
```

`key_file` 文件格式： 每行一个 key, 且仅有 key 的内容，除 key 外不能有其他字符。