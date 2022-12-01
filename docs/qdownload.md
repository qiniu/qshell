# 简介
`qdownload` 可以将七牛空间中的文件同步到本地磁盘中。支持只同步带特定前缀或者后缀的文件，也支持在本地备份路径不变的情况下进行增量同步。 需要额外指出的是，将文件同步到本地都是走的七牛存储源站的流量而不是 `CDN`
的流量，因为 `CDN` 通常情况下会认为大量的文件下载操作是非法访问从而进行限制。

注：
- `Key` 中的 `/` 会被当做路径处理，也即任何以 `/` 结尾的 `Key` 均会被当做文件夹处理。

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
    "save_path_handler"      :   "",
    "prefix"                 :   "image/",
    "suffixes"               :   ".png,.jpg",
    "key_file"               :   "<KeyFile>",
    "check_hash"             :   false,
    "cdn_domain"             :   "down.example.com",
    "referer"                :   "http://www.example.com",
    "public"                 :   true,
    "remove_temp_while_error":   false,
    "log_file"               :   "download.log",
    "log_level"              :   "info",
    "log_rotate"             :   10,
    "log_stdout"             :   false
}
```

字段说明：

- dest_dir：本地数据备份路径，为全路径 【必选】
- bucket：空间名称 【必选】
- prefix：只同步指定前缀的文件，默认为空 【可选】
- suffixes：只同步指定后缀的文件，默认为空 【可选】
- key_file：配置一个文件，指定需要下载的 keys；默认为空，全量下载 bucket 中的文件 【可选】
- save_path_handler：指定一个回调函数；在构建文件的保存路径时，优先使用此选项进行构建，如果不配置则使用 $dest_dir + $文件分割符 + $Key 方式进行构建。文档下面有常用场景实例。此函数通过 Go 语言的模板实现，函数验证使用 func 命令，具体语法可参考 func 命令说明，handler 使用方式下方有示例可供参考 【可选】
- check_hash：是否验证 hash，如果开启可能会耗费较长时间，默认为 `false` 【可选】
- cdn_domain：设置下载的 CDN 域名，默认为空表示从存储源站下载，qshell 下载使用 domain 优先级：1.cdn_domain(此选项) 2.bucket 配置域名(无需配置) 3.qshell 配置文件中 hosts 的 io(需要配置)，当优先级高的 domain 下载失败后会尝试使用优先级低的 domain 进行下载。【该功能默认需要计费，如果希望享受 10G 的免费流量，请自行设置 cdn_domain 参数，如不设置，需支付源站流量费用，无法减免！！！】 【可选】
- referer：如果 CDN 域名配置了域名白名单防盗链，需要指定一个允许访问的 referer 地址；默认为空 【可选】
- public：空间是否为公开空间；为 `true` 时为公有空间，公有空间下载时不会对下载 URL 进行签名，可以提升 CDN 域名性能，默认为 `false`（私有空间）【可选】
- enable_slice: 是否开启切片下载，需要注意 `slice_file_size_threshold` 切片阈值选项的配置，只有开启切片下载，并且下载的文件大小大于切片阈值方会启动切片下载。默认不开启。【可选】
- slice_size: 切片大小；当使用切片下载时，每个切片的大小；单位：B。默认为 4194304，也即 4MB。【可选】
- slice_concurrent_count: 切片下载的并发度；默认为 10 【可选】
- slice_file_size_threshold: 切片下载的文件阈值，当开启切片下载，并且文件大小大于此阈值时方会启用切片下载。【可选】
- remove_temp_while_error: 当下载遇到错误时删除之前下载的部分文件缓存，默认为 `false` (不删除)【可选】
- log_level：下载日志输出级别，可选值为 `debug`,`info`,`warn`,`error`，其他任何字段均会导致不输出日志。默认 `debug` 。【可选】
- log_file：下载日志的输出文件，默认为输出到 `record_root` 指定的文件中，具体文件路径可以在终端输出看到。【可选】
- log_rotate：下载日志文件的切换周期，单位为天，默认为 7 天即切换到新的下载日志文件 【可选】
- log_stdout：下载日志是否同时输出一份到标准终端，默认为 `false`，主要在调试下载功能时可以指定为 `true` 【可选】
- record_root：下载记录信息保存路径，包括日志文件和下载进度文件；默认为 `qshell` 下载目录；【可选】
    - 通过 `-L` 指定工作目录时，`record_root` 则为 `此工作目录/qupload/$jobId`，
    - 未通过 `-L` 指定工作目录时为 `用户目录/.qshell/users/$CurrentUserName/qupload/$jobId`
    - 注意 `jobId` 是根据上传任务动态生成；具体方式为 MD5("DestDir:$Bucket:KeyFile")；`CurrentUserName` 当前用户的名称

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
	"check_hash"    :   false
}
```

运行命令（下载并发数表示可以同时下载 10 个文件）：

```
qshell qdownload -c 10 qdisk_down.conf
```

`key_file` 文件格式： 每行一个 key, 且仅有 key 的内容，除 key 外不能有其他字符。


### `save_path_handler` 说明
`save_path_handler` 函数中可使用的文件参数有：
- Key: 文件的在七牛云存储的 Key 值
- DestDir: 下载配置的保存路径
- ToFile: 默认的下载路径
- ServerFileSize: 文件的在七牛云存储的大小
- ServerFileHash: 文件的在七牛云存储的 Etag
- ServerFilePutTime: 文件的在七牛云存储的上传时间

`save_path_handler` 常见示例：
```
1. 在不配置 save_path_handler 时，文件保存路径的构造方式为：
$dest_dir + $文件分割符 + $Key


2. 配置 save_path_handler ，使在构造文件保存路径时，去除 Key 中一部分前缀 a/：
save_path_handler 配置："{{pathJoin .DestDir (trimPrefix \"a/\" .Key)}}"
pathJoin：路径拼接函数
.DestDir：对应配置文件中的 dest_dir，假设配置为："/user/lala/"
trimPrefix：截掉字符串头函数，trimPrefix \"a/\" .Key 表示：将文件 Key 的 "a/" 截掉
.Key：表示文件的 Key，假设为："a/b/hello.png"
上面信息最终构造的文件路径为："/user/lala/b/hello.png"

如果需要验证 save_path_handler 配置是否符合预期，可使用 func 命令。
参数部分：'{"Key": "a/b/hello.png", "DestDir": "/user/lala/"}'，这部分信息在 download 时会自动生成并作为回调函数的参数，用户不用关心。
回调函数 save_path_handler : "{{pathJoin .DestDir (trimPrefix \"a/\" .Key)}}"
验证：
$qshell func '{"Key": "a/b/hello.png", "DestDir": "/user/lala/"}' "{{pathJoin .DestDir (trimPrefix \"a/\" .Key)}}"
输出：
[W]  output is insert [], and you should be careful with spaces etc.
[I]  [/user/lala/b/hello.png]


3. 自定义文件下载后的保存路径：$DestDir + $文件分割符 + ($Key 首部 a/ 替换成 newA/)
save_path_handler 配置: "{{pathJoin .DestDir \"newA\" (trimPrefix \"a/\" .Key)}}"
pathJoin：路径拼接函数
.DestDir：对应配置文件中的 dest_dir，假设配置为："/user/lala/"
trimPrefix：截掉字符串头函数，trimPrefix \"a/\" .Key 表示：将文件 Key 的 "a/" 截掉
.Key：表示文件的 Key，假设为："a/b/hello.png"
最终文件的保存路径为：
/user/lala/newA/b/hello.png

验证：
$qshell func '{"Key": "a/b/hello.png", "DestDir": "/user/lala/"}' "{{pathJoin .DestDir \"newA\" (trimPrefix \"a/\" .Key)}}"
[W]  output is insert [], and you should be careful with spaces etc.
[I]  [/user/lala/newA/b/hello.png]
```