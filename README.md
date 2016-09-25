# qshell

### 简介
qshell是利用[七牛文档上公开的API](http://d.qiniu.com)实现的一个方便开发者测试和使用七牛API服务的命令行工具。

### 下载

**支持多账户的版本**

**在`1.7.0`及其以上版本，工具所产生的临时文件都存放在工具执行的目录。**

|版本     |支持平台|链接|更新日志|
|--------|---------|----|------|
|qshell v1.8.1|Linux (32, 64位，arm平台), Windows(32, 64位), Mac OSX(32, 64位)|[下载](http://devtools.qiniu.com/qshell-v1.8.1.zip)|[查看](CHANGELOG.md)|


**备注**

`v1.7.0`以上的版本，移除了编译时和平台相关的代码，真正实现代码的跨平台编译，另外从该版本之后`account`命令设置的账号信息将保存在可执行文件执行时的当前目录，而之前的低版本都是在用户目录下面。

**对于`v1.7.0`以上的版本：比如在路径`/Users/jemy/Temp/test`目录下执行命令，那么使用`account`设置的账号信息就在`/Users/jemy/Temp/test`目录下，下次运行命令的时候也需要到这个目录下运行，这种设计的主要目的是方便多账号情况下工具的使用。**

因为上面发布的zip包里面有支持不同平台的可执行文件，请根据系统平台选择合适的可执行文件，然后其他的都可以删除，再把可执行文件重命名为 `qshell` (Windows下面是 `qshell.exe`然后就可以使用了）

###安装

该工具不需要安装，只需要从上面的下载链接下载zip包之后解压即可使用。其中文件名和对应系统关系如下：

|文件名|描述|
|-----|-----|
|qshell_linux_386|Linux 32位系统|
|qshell_linux_amd64|Linux 64位系统|
|qshell_linux_arm|Linux ARM CPU|
|qshell_windows_386.exe|Windows 32位系统|
|qshell_windows_amd64.exe|Windows 64位系统|
|qshell_darwin_386|Mac 32位系统，这种系统很老了|
|qshell_darwin_amd64|Mac 64位系统，主流的系统|

注意，如果在Linux或者Mac系统上遇到`Permission Denied`的错误，请使用命令`chmod +x qshell`来为文件添加可执行权限。这里的`qshell`是上面文件重命名之后的简写。

对于Linux或者Mac，如果希望能够在任何位置都可以执行，那么可以把`qshell`所在的目录加入到环境变量`$PATH`中去。或者最简单的方法如下：

```
sudo mv qshell /usr/local/bin
```
另外，由于本工具是一个命令行工具，在Windows下面请先打开命令行终端，然后输入工具名称执行，不要双击打开。如果你希望可以在任意目录下使用qshell，请将qshell工具可执行文件所在目录添加到系统的环境变量中。

### 使用
我们知道调用七牛的API需要一对`AccessKey`和`SecretKey`，这个可以从七牛的后台的账号设置->[密钥](https://portal.qiniu.com/user/key)获取。

首先要使用七牛的API，必须先设置`AccessKey`和`SecretKey`。命令如下：
```
qshell account ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6o LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKi_
```
上面的`ELUs327kxVPJrGCXqWae9yioc0xYZyrIpbM6Wh6o`就是你的`AccessKey`，而`LVzZY2SqOQ_I_kM1n00ygACVBArDvOWtiLkDtKi_`就是你的`SecretKey`。如果你想查看当前的`AccessKey`和`SecretKey`设置，使用命令：

```
qshell account
```
上面的命令会输出当前你设置好的`AccessKey`和`SecretKey`。
接下来，我们就可以放心地使用七牛的API功能了。

### 选项
可用的选项参数如下：

|参数|描述|
|----|----|
|-f|设置命令行的交互模式，如果指定这个选项，那么为不交互|
|-d|设置是否输出DEBUG日志，如果指定这个选项，则输出DEBUG日志|
|-h|打印命令列表帮助信息|
|-v|打印工具版本|

### 多机房支持
qshell支持通过设置`zone`参数来支持多机房，对应关系如下：

|机房|zone值|
|----|----|
|华东|nb|
|华北|bc|
|华南|hn|
|北美|na0|

### 详解

|命令|描述|详细|
|------|----------|--------|
|account|设置或显示当前用户的AccessKey和SecretKey|[文档](http://github.com/jemygraw/qshell/wiki/account)|
|zone|切换当前设置帐号所在的机房区域，仅账号拥有该指定区域机房时有效|[文档](http://github.com/jemygraw/qshell/wiki/zone)|
|dircache|输出本地指定路径下所有的文件列表|[文档](http://github.com/jemygraw/qshell/wiki/dircache)|
|listbucket|列举七牛空间里面的所有文件|[文档](http://github.com/jemygraw/qshell/wiki/listbucket)|
|alilistbucket|列举阿里OSS空间里面的所有文件|[文档](http://github.com/jemygraw/qshell/wiki/alilistbucket)|
|prefop|查询七牛数据处理的结果|[文档](http://github.com/jemygraw/qshell/wiki/prefop)|
|fput|以文件表单的方式上传一个文件|[文档](http://github.com/jemygraw/qshell/wiki/fput)|
|rput|以分片上传的方式上传一个文件|[文档](http://github.com/jemygraw/qshell/wiki/rput)|
|qupload|同步数据到七牛空间， 带同步进度信息，和数据上传完整性检查|[文档](http://github.com/jemygraw/qshell/wiki/qupload)|
|qdownload|从七牛空间同步数据到本地，支持只同步某些前缀的文件，支持增量同步|[文档](http://github.com/jemygraw/qshell/wiki/qdownload)|
|stat|查询七牛空间中一个文件的基本信息|[文档](http://github.com/jemygraw/qshell/wiki/stat)|
|delete|删除七牛空间中的一个文件|[文档](http://github.com/jemygraw/qshell/wiki/delete)|
|move|移动或重命名七牛空间中的一个文件|[文档](http://github.com/jemygraw/qshell/wiki/move)|
|copy|复制七牛空间中的一个文件|[文档](http://github.com/jemygraw/qshell/wiki/copy)|
|chgm|修改七牛空间中的一个文件的MimeType|[文档](http://github.com/jemygraw/qshell/wiki/chgm)|
|fetch|从Internet上抓取一个资源并存储到七牛空间中|[文档](http://github.com/jemygraw/qshell/wiki/fetch)|
|sync|从Internet上抓取一个资源并存储到七牛空间中，适合大文件的场合|[文档](http://github.com/jemygraw/qshell/wiki/sync)|
|prefetch|更新七牛空间中从源站镜像过来的文件|[文档](http://github.com/jemygraw/qshell/wiki/prefetch)|
|batchdelete|批量删除七牛空间中的文件，可以直接根据`listbucket`的结果来删除|[文档](http://github.com/jemygraw/qshell/wiki/batchdelete)|
|batchchgm|批量修改七牛空间中文件的MimeType|[文档](http://github.com/jemygraw/qshell/wiki/batchchgm)|
|batchcopy|批量复制七牛空间中的文件到另一个空间|[文档](http://github.com/jemygraw/qshell/wiki/batchcopy)|
|batchmove|批量移动七牛空间中的文件到另一个空间|[文档](http://github.com/jemygraw/qshell/wiki/batchmove)|
|batchrename|批量重命名七牛空间中的文件|[文档](http://github.com/jemygraw/qshell/wiki/batchrename)|
|batchrefresh|批量刷新七牛空间中的文件的访问外链|[文档](http://github.com/jemygraw/qshell/wiki/batchrefresh)|
|batchsign|批量根据资源的公开外链生成资源的私有外链|[文档](http://github.com/jemygraw/qshell/wiki/batchsign)|
|checkqrsync|检查qrsync的同步结果，主要通过比对`dircache`和`listbucket`的结果|[文档](http://github.com/jemygraw/qshell/wiki/checkqrsync)|
|b64encode|base64编码工具，可选是否使用UrlSafe方式，默认UrlSafe|[文档](http://github.com/jemygraw/qshell/wiki/b64encode)|
|b64decode|base64解码工具，可选是否使用UrlSafe方式，默认UrlSafe|[文档](http://github.com/jemygraw/qshell/wiki/b64decode)|
|urlencode|url编码工具|[文档](http://github.com/jemygraw/qshell/wiki/urlencode)|
|urldecode|url解码工具|[文档](http://github.com/jemygraw/qshell/wiki/urldecode)|
|ts2d|将timestamp(单位秒)转为UTC+8:00中国日期，主要用来检查上传策略的deadline参数|[文档](http://github.com/jemygraw/qshell/wiki/ts2d)|
|tms2d|将timestamp(单位毫秒)转为UTC+8:00中国日期|[文档](http://github.com/jemygraw/qshell/wiki/tms2d)|
|tns2d|将timestamp(单位100纳秒)转为UTC+8:00中国日期|[文档](http://github.com/jemygraw/qshell/wiki/tns2d)|
|d2ts|将日期转为timestamp(单位秒)|[文档](http://github.com/jemygraw/qshell/wiki/d2ts)|
|ip|根据淘宝的公开API查询ip地址的地理位置|[文档](http://github.com/jemygraw/qshell/wiki/ip)|
|qetag|根据七牛的qetag算法来计算文件的hash|[文档](http://github.com/jemygraw/qshell/wiki/qetag)|
|unzip|解压zip文件，支持UTF-8编码和GBK编码|[文档](http://github.com/jemygraw/qshell/wiki/unzip)|
|privateurl|生成私有空间资源的访问外链|[文档](http://github.com/jemygraw/qshell/wiki/privateurl)|
|saveas|实时处理的saveas链接快捷生成工具|[文档](http://github.com/jemygraw/qshell/wiki/saveas)|
|reqid|七牛自定义头部X-Reqid解码工具|[文档](http://github.com/jemygraw/qshell/wiki/reqid)|
|m3u8delete|根据流媒体播放列表文件删除七牛空间中的流媒体切片|[文档](http://github.com/jemygraw/qshell/wiki/m3u8delete)|
|m3u8replace|修改流媒体播放列表文件中的切片引用域名|[文档](http://github.com/jemygraw/qshell/wiki/m3u8replace)|
|buckets|获取当前账号下所有的空间名称|[文档](http://github.com/jemygraw/qshell/wiki/buckets)|
|domains|获取指定空间的所有关联域名|[文档](http://github.com/jemygraw/qshell/wiki/domains)|
|cdnwho|根据IP地址查询对应的CDN厂商信息，可以用来做域名解析分析|[文档](http://github.com/jemygraw/qshell/wiki/cdnwho)|

##编译
1. 如果是编译本地平台的可执行程序，使用`src`目录下面的`build.sh`脚本即可。
2. 如果是编译跨平台的可执行程序，使用`src`目录下面的`cross_build.sh`脚本即可。

##帮助
如果您遇到任何问题，可以加QQ：2037014430，我将乐意帮助您，非技术问题勿扰。
