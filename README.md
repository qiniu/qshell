# qshell

## 简介
qshell是利用[七牛文档上公开的API](http://developer.qiniu.com)实现的一个方便开发者测试和使用七牛API服务的命令行工具。该工具设计和开发的主要目的就是帮助开发者快速解决问题。目前该工具融合了七牛存储，CDN，以及其他的一些七牛服务中经常使用到的方法对应的便捷命令，比如b64decode，就是用来解码七牛的URL安全的Base64编码用的，所以这是一个面向开发者的工具，任何新的被认为适合加到该工具中的命令需求，都可以在[ISSUE列表](https://github.com/qiniu/qshell/issues)里面提出来，我们会尽快评估实现，以帮助大家更好地使用七牛服务。

## 下载
该工具使用Go语言编写而成，当然为了方便不熟悉Go或者急于使用工具来解决问题的开发者，我们提供了预先编译好的各主流操作系统平台的二进制文件供大家下载使用，由于平台的多样性，我们把这些二进制打包放到一个文件里面，请大家根据下面的说明各自选择合适的版本来使用。在文档中的例子里面，为了方便，我们统一使用`qshell`这个命令来做介绍。

|版本     |支持平台|链接|更新日志|
|--------|---------|----|------|
|qshell v1.8.5|Linux (32, 64位，arm平台), Windows(32, 64位), Mac OSX(32, 64位)|[下载](http://devtools.qiniu.com/qshell-v1.8.5.zip)|[查看](CHANGELOG.md)|

## 安装

该工具由于是命令行工具，所以只需要从上面的下载链接下载zip包之后解压即可使用。其中文件名和对应系统关系如下：

|文件名|描述|
|-----|-----|
|qshell_linux_386|Linux 32位系统|
|qshell_linux_amd64|Linux 64位系统|
|qshell_linux_arm|Linux ARM CPU|
|qshell_windows_386.exe|Windows 32位系统|
|qshell_windows_amd64.exe|Windows 64位系统|
|qshell_darwin_386|Mac 32位系统，这种系统很老了|
|qshell_darwin_amd64|Mac 64位系统，主流的系统|

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

## 使用
该工具有两类命令，其中的一类命令需要指定配置文件，所有的参数信息都写在配置文件里面，比如`qdownload`和`qupload`。还有一类命令的运行需要依赖七牛账号下的`AccessKey`和`SecretKey`，以及空间所在的机房。所以这类命令运行之前，需要使用`account`命令来设置下`AccessKey`，`SecretKey`和默认机房编号`Zone`。

这里需要额外指出的是，为了支持多账号的情况（这种情况很常见，有很多公司有测试账号和线上账号，它们的`AccessKey`和`SecretKey`都不同），我们把`account`设置的参数以及其他命令过程中所产生的临时文件（比如`qupload`会保存已上传文件列表）都存放在`qshell`命令所运行的目录。举个例子，我们有两个目录，一个是测试账号的，一个是线上账号的，分别如下：

```
/home/jemy/tools/qshell_dev/
/home/jemy/tools/qshell_prod/
```

由于上面的特点，我们每次要使用相应账号运行命令的时候，都需要切换到指定的目录下，比如：

```
$ cd /home/jemy/tools/qshell_dev/
$ qshell fput bucket-test test1.png /home/jemy/images/test1.png
```

如果你想看看在命令运行的目录到底保存了哪些文件，除了阅读后面详细的文档外，也可以使用命令`tree -a`查看下。

```
.
├── .qshell
│   ├── account.json
│   └── qupload
│       └── 8bbcb28795215eff7e9362bdd7949b71
│           ├── 8bbcb28795215eff7e9362bdd7949b71.cache
│           ├── 8bbcb28795215eff7e9362bdd7949b71.count
│           └── 8bbcb28795215eff7e9362bdd7949b71.ldb
│               ├── 000001.log
│               ├── CURRENT
│               ├── LOCK
│               ├── LOG
│               └── MANIFEST-000000
```

## 命令选项

该工具还有一些有用的选项参数如下：

|参数|描述|
|----|----|
|-f|设置命令行的交互模式，如果指定这个选项，那么为不交互，这个命令主要用在batch操作的命令里面，因为batch操作具有一定危险性，比如batchdelete批量删除，这个时候会要求输入一个验证码，如果遇到把qshell写到shell脚本里面，不想有验证码这一步的情况，可以指定这个选项|
|-d|设置是否输出DEBUG日志，如果指定这个选项，则输出DEBUG级别的日志|
|-h|打印命令列表帮助信息，遇到参数忘记的情况下，可以使用该命令|
|-v|打印工具版本，反馈问题的时候，请提前告知工具对应版本号|

## 多机房支持

七牛目前上线了华北，华东，华南和北美机房，`qshell`可以同时支持这几个机房的操作，这个是通过`account`命令或者`zone`命令中指定的机房编码来支持的。另外对于`qupload`和`qdownload`，可以在配置文件里面添加选项`zone`来指定空间所在的机房，一旦指定机房之后，工具会自动判定所需要的域名入口信息。目前机房的编码如下：

|机房|zone值|
|----|----|
|华东|nb|
|华北|bc|
|华南|hn|
|北美|na0|

比如通过：

```
$ qshell account ak sk nb 
指定账号默认机房为华东机房

$ qshell zone bc
同一个账号下可能有多个机房的空间，每次都用account比较麻烦，可以用zone命令来切换机房
```

## 命令列表

|命令|类别|描述|详细|
|------|------------|----------|--------|
|account|账号|设置或显示当前用户的`AccessKey`和`SecretKey`和`Zone`|[文档](http://github.com/qiniu/qshell/wiki/account)|
|zone|机房|切换当前设置帐号所在的机房区域，仅账号拥有该指定区域机房时有效|[文档](http://github.com/qiniu/qshell/wiki/zone)|
|dircache|存储|输出本地指定路径下所有的文件列表|[文档](http://github.com/qiniu/qshell/wiki/dircache)|
|listbucket|存储|列举七牛空间里面的所有文件|[文档](http://github.com/qiniu/qshell/wiki/listbucket)|
|prefop|存储|查询七牛数据处理的结果|[文档](http://github.com/qiniu/qshell/wiki/prefop)|
|fput|存储|以文件表单的方式上传一个文件|[文档](http://github.com/qiniu/qshell/wiki/fput)|
|rput|存储|以分片上传的方式上传一个文件|[文档](http://github.com/qiniu/qshell/wiki/rput)|
|qupload|存储|同步数据到七牛空间， 带同步进度信息，和数据上传完整性检查|[文档](http://github.com/qiniu/qshell/wiki/qupload)|
|qdownload|存储|从七牛空间同步数据到本地，支持只同步某些前缀的文件，支持增量同步|[文档](http://github.com/qiniu/qshell/wiki/qdownload)|
|stat|存储|查询七牛空间中一个文件的基本信息|[文档](http://github.com/qiniu/qshell/wiki/stat)|
|delete|存储|删除七牛空间中的一个文件|[文档](http://github.com/qiniu/qshell/wiki/delete)|
|move|存储|移动或重命名七牛空间中的一个文件|[文档](http://github.com/qiniu/qshell/wiki/move)|
|copy|存储|复制七牛空间中的一个文件|[文档](http://github.com/qiniu/qshell/wiki/copy)|
|chgm|存储|修改七牛空间中的一个文件的MimeType|[文档](http://github.com/qiniu/qshell/wiki/chgm)|
|fetch|存储|从Internet上抓取一个资源并存储到七牛空间中|[文档](http://github.com/qiniu/qshell/wiki/fetch)|
|sync|存储|从Internet上抓取一个资源并存储到七牛空间中，适合大文件的场合|[文档](http://github.com/qiniu/qshell/wiki/sync)|
|prefetch|存储|更新七牛空间中从源站镜像过来的文件|[文档](http://github.com/qiniu/qshell/wiki/prefetch)|
|batchdelete|存储|批量删除七牛空间中的文件，可以直接根据`listbucket`的结果来删除|[文档](http://github.com/qiniu/qshell/wiki/batchdelete)|
|batchchgm|存储|批量修改七牛空间中文件的MimeType|[文档](http://github.com/qiniu/qshell/wiki/batchchgm)|
|batchcopy|存储|批量复制七牛空间中的文件到另一个空间|[文档](http://github.com/qiniu/qshell/wiki/batchcopy)|
|batchmove|存储|批量移动七牛空间中的文件到另一个空间|[文档](http://github.com/qiniu/qshell/wiki/batchmove)|
|batchrename|存储|批量重命名七牛空间中的文件|[文档](http://github.com/qiniu/qshell/wiki/batchrename)|
|batchsign|存储|批量根据资源的公开外链生成资源的私有外链|[文档](http://github.com/qiniu/qshell/wiki/batchsign)|
|privateurl|存储|生成私有空间资源的访问外链|[文档](http://github.com/qiniu/qshell/wiki/privateurl)|
|saveas|存储|实时处理的saveas链接快捷生成工具|[文档](http://github.com/qiniu/qshell/wiki/saveas)|
|reqid|存储|七牛自定义头部X-Reqid解码工具|[文档](http://github.com/qiniu/qshell/wiki/reqid)|
|buckets|存储|获取当前账号下所有的空间名称|[文档](http://github.com/qiniu/qshell/wiki/buckets)|
|domains|存储|获取指定空间的所有关联域名|[文档](http://github.com/qiniu/qshell/wiki/domains)|
|qetag|存储|根据七牛的qetag算法来计算文件的hash|[文档](http://github.com/qiniu/qshell/wiki/qetag)|
|m3u8delete|存储|根据流媒体播放列表文件删除七牛空间中的流媒体切片|[文档](http://github.com/qiniu/qshell/wiki/m3u8delete)|
|m3u8replace|存储|修改流媒体播放列表文件中的切片引用域名|[文档](http://github.com/qiniu/qshell/wiki/m3u8replace)|
|cdnrefresh|CDN|批量刷新cdn的访问外链|[文档](http://github.com/qiniu/qshell/wiki/cdnrefresh)|
|cdnprefetch|CDN|批量预取cdn的访问外链|[文档](http://github.com/qiniu/qshell/wiki/cdnprefetch)|
|b64encode|工具|base64编码工具，可选是否使用UrlSafe方式，默认UrlSafe|[文档](http://github.com/qiniu/qshell/wiki/b64encode)|
|b64decode|工具|base64解码工具，可选是否使用UrlSafe方式，默认UrlSafe|[文档](http://github.com/qiniu/qshell/wiki/b64decode)|
|urlencode|工具|url编码工具|[文档](http://github.com/qiniu/qshell/wiki/urlencode)|
|urldecode|工具|url解码工具|[文档](http://github.com/qiniu/qshell/wiki/urldecode)|
|ts2d|工具|将timestamp(单位秒)转为UTC+8:00中国日期，主要用来检查上传策略的deadline参数|[文档](http://github.com/qiniu/qshell/wiki/ts2d)|
|tms2d|工具|将timestamp(单位毫秒)转为UTC+8:00中国日期|[文档](http://github.com/qiniu/qshell/wiki/tms2d)|
|tns2d|工具|将timestamp(单位100纳秒)转为UTC+8:00中国日期|[文档](http://github.com/qiniu/qshell/wiki/tns2d)|
|d2ts|工具|将日期转为timestamp(单位秒)|[文档](http://github.com/qiniu/qshell/wiki/d2ts)|
|ip|工具|根据淘宝的公开API查询ip地址的地理位置|[文档](http://github.com/qiniu/qshell/wiki/ip)|
|unzip|工具|解压zip文件，支持UTF-8编码和GBK编码|[文档](http://github.com/qiniu/qshell/wiki/unzip)|
|alilistbucket|第三方|列举阿里OSS空间里面的所有文件|[文档](http://github.com/qiniu/qshell/wiki/alilistbucket)|

## 项目编译

如果对项目编译感兴趣，请按照如下方式进行：

```
$ go get github.com/syndtr/goleveldb/leveldb
$ go get github.com/yanunon/oss-go-api/oss
$ go get github.com/golang/text
$ ./build.sh
```

## 问题反馈

如果您有任何问题，请写在[ISSUE列表](https://github.com/qiniu/qshell/issues)里面，我们会尽快回复您。

