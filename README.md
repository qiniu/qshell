# qshell

###简介
qshell是利用[七牛文档上公开的API](http://d.qiniu.com)实现的一个方便开发者测试和使用七牛API服务的命令行工具。

###下载

**建议下载最新版本**

|版本     |支持平台|链接|
|--------|---------|----|
|qshell v1.0|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.0.zip)|
|qshell v1.1|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.1.zip)|
|qshell v1.2.1|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.2.1.zip)|
|qshell v1.3|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.3.zip)|
|qshell v1.3.1|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.3.1.zip)|
|qshell v1.3.2|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.3.3.zip)|
|qshell v1.3.3|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.3.3.zip)|
|qshell v1.3.4|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.3.4.zip)|
|qshell v1.3.6|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.3.6.zip)|

###使用
我们知道调用七牛的API需要一对`AccessKey`和`SecretKey`，这个可以从七牛的后台的账号设置->[密钥](https://portal.qiniu.com/setting/key)获取。

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

###详解

|命令|描述|详细|
|------|----------|--------|
|account|设置或显示当前用户的AccessKey和SecretKey|[文档](http://github.com/jemygraw/qshell/wiki/account)|
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
|fetch|从Internet上抓取一个资源到七牛空间中|[文档](http://github.com/jemygraw/qshell/wiki/fetch)|
|prefetch|更新七牛空间中从源站镜像过来的文件|[文档](http://github.com/jemygraw/qshell/wiki/prefetch)|
|batchdelete|批量删除七牛空间中的文件，可以直接根据`listbucket`的结果来删除|[文档](http://github.com/jemygraw/qshell/wiki/batchdelete)|
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

##编译
1. 如果是编译本地平台的可执行程序，使用`src`目录下面的`build.sh`脚本即可。
2. 如果是编译跨平台的可执行程序，使用`src`目录下面的`gox_build.sh`脚本即可。该脚本使用了[gox](https://github.com/mitchellh/gox)工具，请
使用`go get github.com/mitchellh/gox`安装。
