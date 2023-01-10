# 简介
`qshell` 是利用 [七牛文档上公开的API](http://developer.qiniu.com) 实现的一个方便开发者测试和使用七牛 API 服务的命令行工具。该工具设计和开发的主要目的就是帮助开发者快速解决问题。目前该工具融合了七牛存储，CDN，以及其他的一些七牛服务中经常使用到的方法对应的便捷命令，比如 b64decode，就是用来解码七牛的 URL 安全的 Base64 编码用的，所以这是一个面向开发者的工具，任何新的被认为适合加到该工具中的命令需求，都可以在 [ISSUE列表](https://github.com/qiniu/qshell/issues) 里面提出来，我们会尽快评估实现，以帮助大家更好地使用七牛服务。


# 更新日志 
[查看更新文档](https://github.com/qiniu/qshell/blob/master/CHANGELOG.md)


# 命令任何位置运行
**Linux和Mac平台**
对于 Linux 或者 Mac，如果希望能够在任何位置都可以执行，那么可以把 `qshell` 所在的目录加入到环境变量 `$PATH` 中去。假设 `qshell` 命令被解压到路径 `/home/jemy/tools` 目录下面，那么我们可以把如下的命令写入到你所使用的 bash 所对应的配置文件中，如果是 `/bin/bash`，那么就是 `~/.bashrc` 文件，如果是 `/bin/zsh `，那么就是 `~/.zshrc` 文件中。写入的内容为：
```
export PATH=$PATH:/home/jemy/tools
```
保存完毕之后，可以通过两种方式立即生效，其一为输入 `source ~/.zshrc` 或者 `source ~/.bashrc` 来使配置立即生效，或者完全关闭命令行，然后重新打开一个即可，接下来就可以在任何位置使用 `qshell` 命令了。

**Windows平台**
如果你希望可以在任意目录下使用 `qshell`，请将 `qshell` 工具可执行文件所在目录添加到系统的环境变量中。由于 Windows 系统是图形界面，所以方便一点。假设 `qshell.exe` 命令被保存到路径 `E:\jemy\tools` 目录下面，那么我们把这个目录放到系统的环境变量 `PATH` 里面。
注：
本工具是一个命令行工具，在 Windows 下面请先打开命令行终端，然后输入工具名称执行，不要双击打开，否则会出现闪退现象。


# 快速使用
该工具有两类命令，一类需要鉴权，另一类不需要。需要鉴权的命令都需要依赖七牛账号下的 `AccessKey`, `SecretKey`和 `Name`。所以这类命令运行之前，需要使用 `account` 命令来添加 `AccessKey` ，`SecretKey`和`Name` 。
`Name`是用户可以自定义的字符串，用来唯一表示 AccessKey/SecretKey 账户，qshell 会对添加的每一个账户信息加密保存，可以使用自命令 `user` 进行切换，切换账户的时候，需要使用账户唯一标识 `Name`。
```
$ qshell account <Your AccessKey> <Your SecretKey> <Your Name>
```

其中name表示该账号的名称, 如果 ak, sk, name 首字母是 "-", 需要使用如下的方式添加账号, 这样避免把该项识别成命令行选项:
```
$ qshell account -- <Your AccessKey> <Your SecretKey> <Your Name>
```

可以连续使用 qshell account 添加账号 ak, sk, name 信息，qshell 会保存这些账号的信息， 可以使用 qshell user 命令列举账号信息，在各个账号之间切换, 删除账号等。

添加完账户后，就可以使用 qshell 命令对七牛的空间进行操作了。


# 命令选项
该工具还有一些全局的选项，每个命令都可使用，具体如下：
    --doc: 查看命令使用手册，强烈建议使用此选项查看具体命令的文档。eg: qshell user --doc
--silence: 静默模式，控制台仅会输出警告、错误和一些重要的信息
       -d: 设置是否输出 DEBUG 日志，如果指定这个选项，则输出 DEBUG 级别的日志
       -D: 设置是否输出更加详细的 DEBUG 日志，如果指定这个选项，则输出详细的 DEBUG 级别的日志
       -h: 打印命令列表帮助信息，遇到参数忘记的情况下，可以使用该命令
       -v: 打印工具版本，反馈问题的时候，请提前告知工具对应版本号
       -C: qshell配置文件, 其配置格式请看下一节
       -L: 使用当前工作路径作为qshell的配置目录

注：
--silence、-d、-D 优先级：-D > -d > --silence

# 配置文件
1. 配置文件格式支持 json, 如果需要使用配置文件，需要在家目录下创建文件名为 .qshell.json 的 json 文件
2. 配置文件可以配置如 io host, up host, uc host, api host, rs host, rsf host，公有云可以不配置，但私有云必须配置。

例子：
```json
{
    "hosts": {
        "rs": "rs-test.qiniu.com",
        "io": "io-test.qiniu.com",
        "uc": "uc-test.qiniu.com",
        "api": "api-test.qiniu.com",
        "up": "up-test.qiniu.com",
        "rsf": "rsf-test.qiniu.com"
    }
}
```


# 命令列表
### 工具版本信息
- 查看工具版本号请使用命令： qshell -v

### 账号设置命令
- account：设置或显示当前用户的 `AccessKey` 和 `SecretKey`
- user：列举账号信息，在各个账号之间切换，添加账号，删除账号。        

### 存储相关命令
- mkbucket       ：创建存储空间
- bucket         ：查看存储空间信息
- batchdelete    ：批量删除七牛空间中的文件，可以直接根据 `listbucket` 的结果来删除
- delete         ：删除七牛空间中的一个文件
- batchchgm      ：批量修改七牛空间中文件的MimeType
- chgm           ：修改七牛空间中的一个文件的MimeType
- batchchtype    ：批量修改七牛空间中的文件的存储类型
- chtype         ：修改七牛空间中的一个文件的存储类型
- batchexpire    ：批量修改七牛空间中的文件的生存时间
- expire         ：修改七牛空间中的一个文件的生存时间
- batchcopy      ：批量复制七牛空间中的文件到另一个空间
- copy           ：复制七牛空间中的一个文件 
- batchmove      ：批量移动七牛空间中的文件到另一个空间
- move           ：移动或重命名七牛空间中的一个文件
- batchrename    ：批量重命名七牛空间中的文件
- rename         ：重命名七牛空间中的文件
- batchrestorear ：批量解冻七牛空间中的归档/深度归档存储类型文件 
- restorear      ：解冻七牛空间中的归档/深度归档存储类型文件           
- batchstat      ：批量查询七牛空间中文件的基本信息
- stat           ：查询七牛空间中一个文件的基本信息
- buckets        ：获取当前账号下所有的空间名称
- domains        ：获取指定空间的所有关联域名
- listbucket     ：列举七牛空间里面的所有文件 
- listbucket2    ：列举七牛空间里面的所有文件
- batchforbidden ：批量修改文件可访问状态 
- forbidden      ：修改文件可访问状态
- fput           ：以文件表单的方式上传一个文件
- rput           ：以分片上传的方式上传一个文件
- qupload        ：同步数据到七牛空间， 带同步进度信息，和数据上传完整性检查（配置式）
- qupload2       ：同步数据到七牛空间， 带同步进度信息，和数据上传完整性检查（命令式）
- qdownload      ：从七牛空间同步数据到本地，支持只同步某些前缀的文件，支持增量同步 
- get            ：下载存储空间中的文件
- fetch          ：从Internet上抓取一个资源并存储到七牛空间中 
- batchfetch     ：从Internet上抓取一个资源并存储到七牛空间中
- sync           ：从Internet上抓取一个资源并存储到七牛空间中，适合大文件的场合 
- abfetch        ：异步抓取网络资源到七牛存储空间
- m3u8delete     ：根据流媒体播放列表文件删除七牛空间中的流媒体切片 
- m3u8replace    ：修改流媒体播放列表文件中的切片引用域名
- batchsign      ：批量根据资源的公开外链生成资源的私有外链
- dircache       ：输出本地指定路径下所有的文件列表
- prefetch       ：更新七牛空间中从源站镜像过来的文件
- privateurl     ：生成私有空间资源的访问外链 


### CDN 相关的命令
- cdnrefresh  ：批量刷新 cdn 的访问外链或目录
- cdnprefetch ：批量预取 cdn 的访问外链


### 工具类命令
- b64encode ：base64 编码工具，可选是否使用 UrlSafe 方式，默认 UrlSafe
- b64decode ：base64 解码工具，可选是否使用 UrlSafe 方式，默认 UrlSafe
- urlencode ：url 编码工具 
- urldecode ：url 解码工具
- ts2d      ：将 timestamp (单位秒)转为 UTC+8:00 中国日期，主要用来检查上传策略的 deadline 参数
- tms2d     ：将 timestamp (单位毫秒)转为 UTC+8:00 中国日期
- tns2d     ：将 timestamp (单位 100 纳秒)转为 UTC+8:00 中国日期
- d2ts      ：将日期转为 timestamp (单位秒) 
- ip        ：根据淘宝的公开 API 查询 ip 地址的地理位置
- unzip     ：解压 zip 文件，支持 UTF-8 编码和 GBK 编码
- reqid     ：七牛自定义头部 X-Reqid 解码工具
- qetag     ：根据七牛的 qetag 算法来计算文件的 hash
- saveas    ：实时处理的 saveas 链接快捷生成工具
- func      ：封装 Go 语言的模板功能，使用此模板验证 qshell 回调函数逻辑

### 音视频处理相关命令
- pfop   ：提交异步音视频处理请求 
- prefop ：查询七牛数据处理的结果


### 签名类命令
- token ：计算upToken, Qbox token, Qiniu Token


### 其他存储类工具
- alilistbucket ：列举阿里OSS空间里面的所有文件
- awslist       ：列举亚马逊的存储空间
- awsfetch      ：从亚马逊存储迁移数据到七牛存储


# 问题反馈
如果您有任何问题，请写在[ISSUE列表](https://github.com/qiniu/qshell/issues)里面，我们会尽快回复您。
