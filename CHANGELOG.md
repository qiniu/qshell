# 2.6.0
1. 配置文件支持配置 uc host
2. 分片上传 v2 支持指定每个片的大小，具体使用见命令 rput 、qupload 和 qupload2


# 2.5.1
1. 修复在 qnupload 命令上传多个文件时，多次签发 Upload Token 过期时间不变的问题

# 2.5.0
1. 上传功能支持请求 [分片上传 v2](https://developer.qiniu.com/kodo/7458/multipartupload), 默认兼容是原来的 [分片上传 v1](https://developer.qiniu.com/kodo/1650/chunked-upload)，使用姿势见各种上传功能使用文档说明。

# 2.4.3
1. 指定空间为公开空间，[下载](https://github.com/qiniu/qshell/blob/master/docs/qdownload.md)不会对下载 URL 进行签名，可以提升CDN域名性能，默认为私有空间

# 2.4.0
1. 添加batchdelete 自定义分隔符
2. batchfetch支持自定义上传Host设置
3. 添加awsfetch抓取亚马逊空间数据到七牛空间
4. 添加awslist列举亚马逊空间文件
5. 添加了异步抓取命令abfetch

# 2.3.7
1. 加入forbidden命令，可以禁用或者解禁文件
2. pfop加入pipeline, 和上传回调
3. 修复batchfetch bug

# 2.3.6
1. 修复m3u8replace 上传的问题
2. 加入reportCard检测代码质量

# 2.3.5
1. 增加token命令，创建Qbox token, qiniuToken, uploadToken
2. qshell配置文件支持设置accessKey和secretKey
3. 增加了arm版本
4. listbucket2增加文件大小的可读性，可以使用Bytes, KB, MB, GB格式显示

# 2.3.4
1. listbucket2增加捕捉interrupt信号(CTR-C), 打印marker
2. account在本地记录账号，默认不覆盖, 加了-w强制覆盖选项
3. listbucket2 增加append 模式（-a)开启, 修复列举几亿空间的时候，列举一半左右程序中断问题
4. 修复dircache 列表没有输出到文件使用-o选项的时候
5. 修复qupload, qupload2使用多线程上传导致的部分文件上传失败问题
6. 加了-L 选项到qshell, 使用当前工作路径作为qshell的配置目录

# 2.3.3
1. 修复qdownload配置cdn_domain使用了测试域名作为HOST 引起超过10G流量限制的问题
2. listbucket2 max-retry选项只限制出错下载次数，不限制接口返回空的次数
3. 修复cdnprefetch, cdnrefresh中的bug
4. 增加rput, fput, qupload的设置callbackurl, callbackhost的支持

# 2.3.2
1. 修复fetch incorrect region
2. 修复docs中文档格式显示问题
3. 给listbucket2添加不限制重试次数的功能

# 2.3.1
1. batchdelete, batchchgm, batchchtype, batchmove, batchrename命令运行可导出失败，成功文件列表
2. rput, fput, qupload上传支持设置回调用
3. 修复qdownload, qupload 配置文件Windows下UTF-8 BOM解码出错问题

# 2.3.0
1. 重构qshell代码，方便后续合并qrsctl, qfetch，和添加新的功能
2. 加入了user命令， 多用户管理，可以在多个用户之间切换
3. 加入了命令pfop, 命令行提交音视频的处理
4. get 命令，这个是移植qrsctl get的命令的，下载文件
5. 增加了包依赖管理，（go1.11以上）
6. 加入了命令行自动补全， 用户切换时候用户名字的补全
7. 增加了默认读取家目录下.qshell.json格式的配置文件，这样对于通用的配置，可以不用每次都要命令行提供配置文件
8. 修改了命令行接口行为, 选项是以"-"开头的后跟一个字母， 长选项--开头, 对于没有提供配置文件的命令行，默认从标准输入读取内容

# 2.2.1
1. 为qdownload增加file_encoding参数用来支持windows下的gbk编码
2. 为qupload增加file_encoding参数用来支持windows下的gbk编码

# 2.2.0
1. 增加get命令直接从源站下载文件
2. 增加listbucket2命令支持流式获取空间文件列表

# 2.1.9
1. 增加获取授权空间域名的功能

# 2.1.8
1. 增加新加坡机房的支持
2. 全平台使用原生平台支持

# 2.1.7
1. 增加`rpcdecode`功能，用来解码qiniu编码参数，支持按行解码
2. 增加`rpcencode`功能，用来支持qiniu编码
3. 修复`qdownload`对私有云和公有云的兼容性
4. 增加`-c`选项，支持从指定的账户文件读取信息

# 2.1.6
1. 为资源管理操作添加keepalive的选项，支持海量文件快速管理
2. 为上传开启keepalive的选项，优化上传的性能

# 2.1.5
1. 增加私有云支持，可以通过`-f`选项指定host文件
2. 修复`qdownload`下载文件名中有特殊字符导致404的问题

# 2.1.4
1. 优化`sync`命令检测服务端是否支持Range的算法

# 2.1.3
1. 修复私有链接签名的时候文件名存在空格的时候的问题
2. 增加`chtype`命令，支持修改文件的存储类型
3. 增加`batchchtype`命令，支持批量修改文件的存储类型
4. 增加`expire`命令，支持修改文件的生存时间
5. 增加`batchexpire`命令，支持批量修改文件的生存时间

# 2.1.2
1. 优化`qupload`的导出的文件列表的写入文件的及时性

# 2.1.1
1. 为`qupload`增加`delete_on_success`配置，支持上传成功后删除本地文件
2. 为`qupload`增加`success-list`，`failure-list`以及`overwrite-list`几个命令选项，支持导出上传结果的文件列表
3. 调大分片上传的时候并发的块的数量，支持大磁盘IO和大网络IO的场景

# 2.1.0
1. 修复`up_host`指定的上传域名存在尾部斜杠的时候造成的上传错误
2. 为表单上传方式默认设置`crc32`检测，以避免内容损坏

# 2.0.9
1. 修复`sync`命令发送最后一块`Range`请求指定的头部不严格的问题
2. 修复`alilistbucket`无法列举OSS空间列表的问题

# 2.0.8
1. `fput`, `rput`, `stat`, `qupload`, `qupload2`，`batchstat`，`listbucket`命令支持低频存储
2. 修复`sync`命令不支持链接302的bug

# 2.0.7
1. 修复`qdownload`使用`cdn_domain`下载时，可能出现的400和404错误。

# 2.0.6
1. 为`batchcopy`，`batchdelete`，`batchmove`，`batchchgm`，`batchrename`添加并发操作参数
2. 为`qdownload`添加`cdn_domain`和`referer`配置参数，允许用户从CDN域名下载资源

# 2.0.5
1. 修复2.0.1引入的bug，该bug不影响文件上传结果，但是影响文件上传效率。主要是
决定采用表单上传还是分片上传的时候逻辑写反了。

# 2.0.4
1. 修复windows下面beego的log报错的问题

# 2.0.3
1. 修改各服务的默认域名，使用机房相关的域名

# 2.0.2
1. 为命令添加exit code，可以使用 `echo $?` 来查看命令执行结果

# 2.0.1
1. 修复一个`dircache`可能导致runtime错误的bug
2. 使用本地环境编译linux 386和amd64平台的二进制文件

# 2.0.0
1. 这是一个全新设计过的版本，支持单用户模式和多用户模式
2. 抛弃了手动设置空间机房的方式，直接根据空间名称自动获取对应机房
3. 所有的上传操作默认走空间所在机房的上传加速域名
4. `qdownload`目前只支持从源站批量下载文件，并支持断点续传
5. 将文档从wiki迁移到项目目录`docs`下，方便随时修改

# 1.8.5
1. 更新rs操作遇到错误的时候的输出，更加好看一点
2. 添加`cdnrefresh`和`cdnprefetch`功能

# 1.8.4
1. 修复 rs 域名问题

#1.8.3
1. 优化`SetZone`，方便未来扩展到更多机房
2. 修复`zone`指令设置`hn`没有生效的问题

#1.8.2
1. 增加对华南机房的支持，`zone`和`account`指令支持新zone为`hn`

#1.8.1
1. 增加对北美机房的支持，`zone`和`account`指令支持新zone为`na0`
2. 优化`qdownload`功能，后面一律采用代理模式走源站下载文件

#1.8.0
1. 更新`qupload`和`qdownload`的并发模型，支持更高效的并发上传和下载

#1.7.9
1. 为`qdownload`添加`zone`参数，支持多机房文件下载，可选值为`nb`，`bc`，`aws`，默认为`nb`

#1.7.8
1. 添加`m3u8replace`命令，可以修改m3u8文件中的域名信息

#1.7.7
1. 添加`-f`参数，可以在使用批量处理功能的时候，不必输入确认码

#1.7.6
1. 添加`cdnwho`指令，可以根据IP地址查询出对应的CDN厂商信息
2. 添加`regionip`指令，可以根据地区和运营商查询CDN域名解析的IP列表

#1.7.5
1. 为`qupload`添加`check_hash`参数，默认在`check_exists`情况下不进行hash匹配，节约同步时间
2. 为`qupload`添加`check_size`参数，默认在`check_exists`情况下如果不进行hash匹配，可以启用文件大小匹配，做简单匹配，节约同步时间的同时，检测文件是否变化
3. 为`qupload`添加`log_level`参数，可以设置日志级别，可以在同步出错的情况下，只输出错误信息方便查看原因
4. 为`qupload`添加`log_file`参数，默认输出到终端，如果设置，则输出日志到指定的本地磁盘文件中
5. 为`qupload`添加`skip_fixed_strings`参数，可以根据指定的字符串忽略相对路径中存在这些字符串的文件，不进行同步

# 1.7.4
1. 为`qdownload`添加`referer`参数，支持防盗链开启情况下的资源下载

#1.7.3
1. 使`prefop`支持不同的zone
2. 更新`rput`和`qupload`的上传token有效期为7天

#1.7.2
1. 为`qupload`添加配置参数`file_list`，可以从文件中读取待同步列表

#1.7.1
1. 移除项目编译时和平台相关的代码，真正实现跨平台编译。
2. 支持Linux ARM CPU。
3. 使得`account`命令和`qupload`生成的临时文件保存在工具当前执行目录之下
4. 添加`qupload2`命令，功能和`qupload`相同，只是参数通过命令行参数来指定，方面在shell脚本中使用

#1.6.5
1. 修复32位系统下面，Unix时间戳存储和解析时被截断造成的文件覆盖错误

#1.6.4
1. 优化`fput`的上传进度输出
2. 优化资源管理相关命令的错误输出
3. 优化`prefop`的输出
4. 移除断点续传中的冗余的逻辑

#1.6.3
1. 修复`m3u8delete`命令的请求域名bug

#1.6.2
1. 修复`account`命令设置`zone`无法生效的bug
2. 优化`rput`和`fput`的进度显示和错误提示信息
3. 为`qupload`添加断点续传功能，对大文件的上传比较友善
4. 修复`qupload`分片上传时，不会记录已上传的文件导致再次同步时有可能重复上传

#1.6.1
1. 为`qupload`细分根据相对路径还是文件名来跳过文件不上传

#1.6.0
1. 修复`sync`同步大文件时，token过期的问题

#1.5.9
1. 添加命令`sync`，支持将大文件通过Range的方式迁移到七牛

#1.5.8
1. 为`rput`和`fput`添加覆盖上传功能

#1.5.7
1. 分片上传片大小调整为`4M`，以便于最大化利用带宽

#1.5.6
1. 为工具添加`zone`功能，可以支持多机房的操作
2. 将七牛的golang sdk v6下载到本地引入，以方便后面修改
3. 为`qupload`添加本地网卡绑定功能`bind_nic_ip`
4. 为`qupload`添加远程主机IP绑定功能`bind_up_ip`和`bind_rs_ip`，避免频繁DNS
5. 为`qupload`添加选项参数`put_threshold`，用来指定分片上传的阈值
6. 添加命令`zone`用来输出当前机房信息或者切换机房

#1.5.5
1. 修复`dir_cache`中的`WalkFunc`没有处理文件遍历错误的bug
2. 为`qupload`添加文件遍历提示信息，在大批量文件遍历的情况下，比较友好
3. 修复`unzip`中的bug，支持pkzip 4.5格式解压

#1.5.4
1. 改进文件路径拼接代码，提高跨平台友好性

#1.5.3
1. 添加批量签名的功能，可以在迁移私有空间资源的时候使用，批量签发访问链接

#1.5.2
1. 添加`batchrefresh`指令，可以用来根据资源外链批量刷新cdn缓存

#1.5.1
1. 修复`qupload`中的一个严重bug，该bug会导致本地目录下如果有一个文件无法访问而导致的上传中断

#1.5.0
1. 优化m3u8获取域名的逻辑
2. Batch操作的验证在Windows的操作系统下面不显示颜色

#1.4.9
1. 为`batchdelete`,`batchcopy`,`batchrename`,`batchchgm`,`batchmove`添加操作确认功能

#1.4.8
1. 修复`qupload`的`src_dir`参数指定尾部`/`或者`\\`造成本地文件名多余问题
2. 修复`qupload`的`check_exists`在远程文件删除，本地已同步过的情况下，再次同步不会上传的问题
3. 为`qupload`添加`skip_prefixes`和`skip_suffixes`配置参数，用来支持根据文件或目录的前缀或后缀忽略文件或目录，不上传
4. 使用系统文件路径拼接方法，替换原有的手动字符串拼接
5. 修复日志的输出重定向无法到文件中的问题

#1.4.7
1. 为`qdownload`添加`suffix`选项，可以进一步缩小需要下载的文件列表，只有以`suffix`结尾的才会被下载
2. 修改`qupload`的分片上传的阈值为10MB，这样有利于提高服务器端上传的有效性，尤其是海外上传
3. 修改`qupload`的分片大小为1MB，以提升低带宽情况下，上传的有效性，尤其是海外上传
4. 修改`qupload`的<ThreadCount>提示信息的日志级别，从WARN换成INFO，避免给客户造成困扰
5. 更新七牛go sdk的引用路径，使用v6的版本编译

#1.4.6
1. 修复`qupload`命令在同步路径太长情况下造成本地日志文件无法打开的错误，使用md5的方式生成新的日志名称
2. 添加指令`buckets`来获取当前账号下的所有空间名称
3. 添加指令`domains`来获取指定`bucket`下面的所有相关域名

#1.4.5
1. 修复`qupload`命令`check_exists`参数的判断逻辑错误。

#1.4.4
1. 为`rput`和`fput`添加可选参数`UpHost`用来指定特定的上传入口

#1.4.3
1. 为`qupload`添加`check_exists`参数，该参数为true的时候，会在每个文件进行上传的时候，发送stat请求来查询七牛空间是否已经有这个文件，如果没有才上传。
2. 添加`-h`和`-v`选项，让`qshell`看上去和其他的指令一样
3. 修正了`qupload`进度的格式化问题

#1.4.2
1. 修复`m3u8delete`指令在私有空间下生成资源外链的错误
2. 对于`batchdelete`,`batchchgm`,`batchcopy`,`batchmove`,`batchrename`的配置文件，不再过滤字段两端的空白字符，以支持某些特殊的文件名称
3. 使上传入口可配置，添加了配置选项`up_host`，用户可以选择性设置不同的上传入口来满足他们不同的需求，比如海外文件同步等

#1.4.1
1. 添加指令`m3u8delete`来支持在删除`m3u8`文件的同时，删除所引用的切片文件
2. 优化`prefop`的内容输出

#1.4.0
1. 添加指令`batchcopy`支持云端空间内或者空间到空间的资源批量复制功能
2. 添加指令`batchcopy`的文档说明
1. 添加指令`reqid`来解析七牛http回复的头部字段`X-Reqid`
2. 改进代码，修复windows平台文件路径分隔符的问题
3. 改进代码，自动删除`qupload`指令的旧的本地磁盘列表缓存文件
4. 删除旧的发布版本，直接使用新的发布版本，以保证和文档的一致性
5. 添加`batchchgm`,`batchmove`,`batchrename`,`saveas`,`reqid`的文档说明

#1.3.9
1. 添加指令`batchmove`支持云端空间内或者空间到空间的资源批量移动功能
2. 添加指令`batchrename`支持云端空间内或者空间到空间的资源批量重命名功能
3. 添加指令`batchchgm`支持云端空间内文件的MimeType批量修改功能
4. 添加指令`saveas`,可以用来快速生成以外链方式进行的文件处理并保存的链接
5. 更新`fetch`功能以保持和最新的API一致，最新的`fetch`API返回内容是文件的hash和key

#1.3.8
1. 将指令`privateurl`的参数`Deadline`变为可选参数，默认距当前时间1小时后的时间戳
2. 为指令`qupload`添加上传前文件的最后修改时间的检测，如果配置开启了`overwrite`选项，那么在文件的最后修改时间和上次不同时，将覆盖上传
3. 移除`qupload`同步完成的文件上传完整性自动检测功能，这个功能在海量文件同步情况下会很慢
4. 移除`unzip`原有的对`iconv`的编译依赖，使得其支持windows的平台
5. 使用`SrcDir+":"+Bucket`来构成一个默认的上传任务id，使得上传的检测更加灵活

#1.3.7
1. 添加指令`qetag`用来根据七牛的文件hash算法来计算一个本地文件的hash
2. 为指令`qupload`添加并发参数支持，可选范围1-100
3. 添加指令`unzip`来支持解压七牛生成的zip文件，支持gbk和utf8编码
4. 添加指令`privateurl`用来为七牛私有空间的资源快速生成一个外链
5. 添加`qetag`,`unzip`和`privateurl`的文档说明

#1.3.6
1. 为指令`qupload`添加参数`overwrite`用来支持本地文件改动后的覆盖上传
2. 修复`batchdelete`的bug，该bug导致文件数量超过1000时，只删除1000个文件

#1.3.4
1. 为指令`fput`和`rput`添加上传过程中平均速度的功能
2. 优化指令`fput`和`rput`上传过程平均速度的显示

#1.3.3
1. 修复windows平台下文件命令的路径分隔符问题
2. 为指令`prefop`的输出结果添加`Bucket`字段

#1.3.2
1. 修复windows下文件路径分隔符的bug
2. 修复指令`qdownload`的帮助信息bug

#1.3
1. 添加指令`qdownload`用来从七牛云端批量同步文件到本地
2. 修复`qdownload`在私有空间的情况下，生成token的bug
3. 检查`qdownload`下载文件的时候，每个请求的状态码是否为200

#1.2.1
1. 为readme文档添加项目编译信息
2. 修复bufio的Writer文件内容没有自动写入的bug，需要Flush
3. 为指令`prefop`添加详细的输出信息

#1.2
1. 将指令`b64encode`和`b64decode`的参数`UrlSafe`变为可选，默认为true
2. 将`qupload`的日志信息按照任务不同，独立到单独的文件夹中
3. 添加指令`checkqrsync`用来检查`qrsync`同步工具的数据同步完整性

#1.1
1. 添加指令`account`用来设置`AccessKey`和`SecretKey`
2. 添加指令`dircache`用来获取本地目录下的文件列表
3. 添加指令`listbucket`用来从七牛的空间中获取资源的列表
4. 添加指令`prefop`用来查询七牛fop数据处理结果
5. 添加指令`stat`用来获取七牛空间中文件的基本信息
6. 添加指令`delete`用来从七牛空间删除文件
7. 添加指令`move`用来在空间内或者空间和空间之间移动文件
8. 添加指令`copy`用来在空间内或者空间和空间之间复制文件
9. 添加指令`chgm`用来改变空间中文件的MimeType
10. 添加指令`fetch`用来支持将文件抓取到空间中
11. 添加指令`prefetch`用来支持从镜像源站更新七牛空间中的文件
12. 添加指令`batchdelete`用来从七牛空间批量删除文件
13. 添加指令`fput`用来支持以表单方式上传文件到空间
14. 添加指令`rput`用来支持以分片方式上传文件到空间
15. 添加指令`alilistbucket`用来从阿里云OSS空间中获取资源的列表
16. 添加指令`b64encode`用来支持对字符串进行base64或者url安全的base64编码
17. 添加指令`b64decode`用来支持对字符串进行base64或者url安全的base64解码
18. 添加指令`urlencode`用来支持对字符串进行url编码
19. 添加指令`urldecode`用来支持对字符串进行url解码
20. 添加指令`ts2d`用来将以秒为单位的时间戳转化为可读的日期
21. 添加指令`tns2d`用来将以100纳秒为单位的时间戳转化为可读日期
22. 添加指令`d2ts`用来将当前的本地日志转化为以秒为单位的Unix时间戳
23. 添加指令`ip`用来支持查询一个或者多个ip地址的供应商信息
24. 添加指令`qupload`用来支持将本地目录下的文件批量同步到七牛空间
