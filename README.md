# qshell

## 简介

qshell是利用[七牛文档上公开的API](http://developer.qiniu.com)实现的一个方便开发者测试和使用七牛API服务的命令行工具。该工具设计和开发的主要目的就是帮助开发者快速解决问题。目前该工具融合了七牛存储，CDN，以及其他的一些七牛服务中经常使用到的方法对应的便捷命令，比如b64decode，就是用来解码七牛的URL安全的Base64编码用的，所以这是一个面向开发者的工具，任何新的被认为适合加到该工具中的命令需求，都可以在[ISSUE列表](https://github.com/qiniu/qshell/issues)里面提出来，我们会尽快评估实现，以帮助大家更好地使用七牛服务。

## 下载

该工具使用Go语言编写而成，当然为了方便不熟悉Go或者急于使用工具来解决问题的开发者，我们提供了预先编译好的各主流操作系统平台的二进制文件供大家下载使用，由于平台的多样性，我们把这些二进制打包放到一个文件里面，请大家根据下面的说明各自选择合适的版本来使用。在文档中的例子里面，为了方便，我们统一使用`qshell`这个命令来做介绍。

> 更新日志 [查看](CHANGELOG.md)

|版本     |支持平台|链接|
|--------|---------|----|
|qshell v2.3.0|Mac OSX, Linux, Windows|[下载](http://devtools.qiniu.com/qshell-v2.3.0.zip)|

## 安装

该工具由于是命令行工具，所以只需要从上面的下载之后即可使用。其中文件名和对应系统关系如下：

|文件名|描述|
|-----|-----|
|qshell_linux_x86 |Linux 32位系统|
|qshell_linux_x64|Linux 64位系统|
|qshell_windows_x86.exe|Windows 32位系统|
|qshell_windows_x64.exe|Windows 64位系统|
|qshell_darwin_x64|Mac 64位系统，主流的系统|

**Linux和Mac平台**

（1）权限
如果在Linux或者Mac系统上遇到`Permission Denied`的错误，请使用命令`chmod +x qshell`来为文件添加可执行权限。这里的`qshell`是上面文件重命名之后的简写。

（2）任何位置运行
对于Linux或者Mac，如果希望能够在任何位置都可以执行，那么可以把`qshell`所在的目录加入到环境变量`$PATH`中去。假设`qshell`命令被解压到路径`/home/jemy/tools`目录下面，那么我们可以把如下的命令写入到你所使用的bash所对应的配置文件中，如果是`/bin/bash`，那么就是`~/.bashrc`文件，如果是`/bin/zsh`，那么就是`~/.zshrc`文件中。写入的内容为：

```
export PATH=$PATH:/home/jemy/tools
```
保存完毕之后，可以通过两种方式立即生效，其一为输入`source ~/.zshrc`或者`source ~/.bashrc`来使配置立即生效，或者完全关闭命令行，然后重新打开一个即可，接下来就可以在任何位置使用`qshell`命令了。

**Windows平台**

（1）闪退问题
本工具是一个命令行工具，在Windows下面请先打开命令行终端，然后输入工具名称执行，不要双击打开，否则会出现闪退现象。

（2）任何位置运行
如果你希望可以在任意目录下使用`qshell`，请将`qshell`工具可执行文件所在目录添加到系统的环境变量中。由于Windows系统是图形界面，所以方便一点。假设`qshell.exe`命令被解压到路径`E:\jemy\tools`目录下面，那么我们把这个目录放到系统的环境变量`PATH`里面。

![windows-qshell-path-settings.png](http://devtools.qiniu.com/windows-qshell-path-settings.png)

## qshell使用

1. 添加密钥和账户名称

该工具有两类命令，一类需要鉴权，另一类不需要。

需要鉴权的命令都需要依赖七牛账号下的 `AccessKey`, `SecretKey`和 `Name`。所以这类命令运行之前，需要使用 `account` 命令来添加 `AccessKey` ，`SecretKey`和`Name` 。
`Name`是用户可以自定义的字符串，用来唯一表示AccessKey/SecretKey账户，qshell会对添加的每一个账户信息加密保存，可以使用自命令`user`进行切换，切换账户的时候，需要使用账户唯一标识
`Name`。 

```
$ qshell account <Your AccessKey> <Your SecretKey> <Your Name>
```

其中name表示该账号的名称, 如果ak, sk, name首字母是"-", 需要使用如下的方式添加账号, 这样避免把该项识别成命令行选项:

```
$ qshell account -- <Your AccessKey> <Your SecretKey> <Your Name>
```

可以连续使用qshell account 添加账号ak, sk, name信息，qshell会保存这些账号的信息， 可以使用qshell user命令列举账号信息，在各个账号之间切换, 删除账号等

2. 添加完账户后，就可以使用qshell上传，下载文件了

## 账户管理

使用qshell user子命令可以用来管理记录的多账户信息。
1. qshell user ls可以列举账户下所有的账户信息
2. qshell user cu <userName>可以用来切换账户
3. qshell user cu 不携带<userName>的话会切换到最近的上个账户；比如我在A账户做完操作后，使用qshell user cu B到了B 账户，那么使用qshell user cu可以切回到A账户.

## 开启命令的自动补全
**linux上，使用bash**

在centos上，需要安装bash-completion包，默认该包没有安装

```
$ yum install bash-completion -y

```

给当前的bash加入自动补全， 运行命令

```
$ source <(qshell completion bash)

```

也可以把这个配置加入bash的启动文件中，这样到bash启动的时候，会自动加载qshell的补全配置：


```
$ echo "source <(qshell completion bash)" >> ~/.bashrc
```

**Mac上，使用zsh**
把如下代码加入zsh的启动文件中~/.zshrc, 然后source ~/.zshrc

```
if [ $commands[qshell] ]; then
  source <(qshell completion zsh)
fi
```

## 命令选项

该工具还有一些有用的选项参数如下：

|参数|描述|
|----|----|
|-d|设置是否输出DEBUG日志，如果指定这个选项，则输出DEBUG级别的日志|
|-h|打印命令列表帮助信息，遇到参数忘记的情况下，可以使用该命令|
|-v|打印工具版本，反馈问题的时候，请提前告知工具对应版本号|
|-C|qshell配置文件, 其配置格式请看下一节|

## 配置文件

1. 配置文件格式支持json, 如果需要使用配置文件，需要在家目录下创建文件名为.qshell.json的json文件
2. 配置文件可以配置如io host, api host, rs hsot, rsf host, 这些如果没有指定，程序会自动选择

例子：

默认官方的列举空间的文件使用的是rs.qiniu.com域名，如果因为某种原因，比如私有存储，需要替换使用rs-test.qiniu.com这个域名的话，那么只需要
在家目录下创建文件名字为.qshell.json的配置文件，文件内容为

```json
{
    "hosts": {
        "rs_host": "rs-test.qiniu.com"
    }
}
```

如果想要更改io host为io-test.qiniu.com的话，只需要继续在上面的hosts中添加，如下：

```json
{
    "hosts": {
        "rs_host": "rs-test.qiniu.com",
        "io_host": "io-test.qiniu.com"
    }
}
```

同理如果全部修改的话

```json
{
    "hosts": {
        "rs_host": "rs-test.qiniu.com",
        "io_host": "io-test.qiniu.com",
        "api_host": "",
        "rsf_host": ""
    }
}
```


## 命令列表

|命令|类别|描述|详细|
|------|------------|----------|--------|
|account|账号|设置或显示当前用户的`AccessKey`和`SecretKey`|[文档](docs/account.md)|
|dircache|存储|输出本地指定路径下所有的文件列表|[文档](docs/dircache.md)|
|listbucket|存储|列举七牛空间里面的所有文件|[文档](docs/listbucket.md)|
|listbucket2|存储|列举七牛空间里面的所有文件|[文档](docs/listbucket2.md)|
|prefop|存储|查询七牛数据处理的结果|[文档](docs/prefop.md)|
|fput|存储|以文件表单的方式上传一个文件|[文档](docs/fput.md)|
|rput|存储|以分片上传的方式上传一个文件|[文档](docs/rput.md)|
|qupload|存储|同步数据到七牛空间， 带同步进度信息，和数据上传完整性检查（配置式）|[文档](docs/qupload.md)|
|qupload2|存储|同步数据到七牛空间， 带同步进度信息，和数据上传完整性检查（命令式）|[文档](docs/qupload2.md)|
|qdownload|存储|从七牛空间同步数据到本地，支持只同步某些前缀的文件，支持增量同步|[文档](docs/qdownload.md)|
|stat|存储|查询七牛空间中一个文件的基本信息|[文档](docs/stat.md)|
|delete|存储|删除七牛空间中的一个文件|[文档](docs/delete.md)|
|move|存储|移动或重命名七牛空间中的一个文件|[文档](docs/move.md)|
|copy|存储|复制七牛空间中的一个文件|[文档](docs/copy.md)|
|chgm|存储|修改七牛空间中的一个文件的MimeType|[文档](docs/chgm.md)|
|chtype|存储|修改七牛空间中的一个文件的存储类型，支持普通存储（0）和低频存储（1）|[文档](docs/chtype.md)|
|expire|存储|修改七牛空间中的一个文件的生存时间|[文档](docs/expire.md)|
|fetch|存储|从Internet上抓取一个资源并存储到七牛空间中|[文档](docs/fetch.md)|
|sync|存储|从Internet上抓取一个资源并存储到七牛空间中，适合大文件的场合|[文档](docs/sync.md)|
|prefetch|存储|更新七牛空间中从源站镜像过来的文件|[文档](docs/prefetch.md)|
|batchdelete|存储|批量删除七牛空间中的文件，可以直接根据`listbucket`的结果来删除|[文档](docs/batchdelete.md)|
|batchchgm|存储|批量修改七牛空间中文件的MimeType|[文档](docs/batchchgm.md)|
|batchchtype|存储|批量修改七牛空间中的文件的存储类型，支持普通存储（0）和低频存储（1）|[文档](docs/batchchtype.md)|
|batchexpire|存储|批量修改七牛空间中的文件的生存时间|[文档](docs/batchexpire.md)|
|batchcopy|存储|批量复制七牛空间中的文件到另一个空间|[文档](docs/batchcopy.md)|
|batchmove|存储|批量移动七牛空间中的文件到另一个空间|[文档](docs/batchmove.md)|
|batchrename|存储|批量重命名七牛空间中的文件|[文档](docs/batchrename.md)|
|batchsign|存储|批量根据资源的公开外链生成资源的私有外链|[文档](docs/batchsign.md)|
|batchstat|存储|批量查询七牛空间中文件的基本信息|[文档](docs/batchstat.md)|
|privateurl|存储|生成私有空间资源的访问外链|[文档](docs/privateurl.md)|
|saveas|存储|实时处理的saveas链接快捷生成工具|[文档](docs/saveas.md)|
|reqid|存储|七牛自定义头部X-Reqid解码工具|[文档](docs/reqid.md)|
|buckets|存储|获取当前账号下所有的空间名称|[文档](docs/buckets.md)|
|domains|存储|获取指定空间的所有关联域名|[文档](docs/domains.md)|
|qetag|存储|根据七牛的qetag算法来计算文件的hash|[文档](docs/qetag.md)|
|m3u8delete|存储|根据流媒体播放列表文件删除七牛空间中的流媒体切片|[文档](docs/m3u8delete.md)|
|m3u8replace|存储|修改流媒体播放列表文件中的切片引用域名|[文档](docs/m3u8replace.md)|
|cdnrefresh|CDN|批量刷新cdn的访问外链或目录|[文档](docs/cdnrefresh.md)|
|cdnprefetch|CDN|批量预取cdn的访问外链|[文档](docs/cdnprefetch.md)|
|b64encode|工具|base64编码工具，可选是否使用UrlSafe方式，默认UrlSafe|[文档](docs/b64encode.md)|
|b64decode|工具|base64解码工具，可选是否使用UrlSafe方式，默认UrlSafe|[文档](docs/b64decode.md)|
|urlencode|工具|url编码工具|[文档](docs/urlencode.md)|
|urldecode|工具|url解码工具|[文档](docs/urldecode.md)|
|ts2d|工具|将timestamp(单位秒)转为UTC+8:00中国日期，主要用来检查上传策略的deadline参数|[文档](docs/ts2d.md)|
|tms2d|工具|将timestamp(单位毫秒)转为UTC+8:00中国日期|[文档](docs/tms2d.md)|
|tns2d|工具|将timestamp(单位100纳秒)转为UTC+8:00中国日期|[文档](docs/tns2d.md)|
|d2ts|工具|将日期转为timestamp(单位秒)|[文档](docs/d2ts.md)|
|ip|工具|根据淘宝的公开API查询ip地址的地理位置|[文档](docs/ip.md)|
|unzip|工具|解压zip文件，支持UTF-8编码和GBK编码|[文档](docs/unzip.md)|
|alilistbucket|第三方|列举阿里OSS空间里面的所有文件|[文档](docs/alilistbucket.md)|
|get|存储|下载存储空间中的文件|[文档](docs/get.md)|
|pfop|dora异步处理|提交异步音视频处理请求|[文档](docs/pfop.md)

## 问题反馈

如果您有任何问题，请写在[ISSUE列表](https://github.com/qiniu/qshell/issues)里面，我们会尽快回复您。
