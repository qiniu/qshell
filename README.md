# qshell

###简介
qshell是利用[七牛文档上公开的API](http://d.qiniu.com)实现的一个方便开发者测试和使用七牛API服务的命令行工具。

###下载
|版本     |支持平台|链接|
|--------|---------|----|
|qshell v1.0|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qshell v1.0.zip)|

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
**dircache - 获取本地系统指定路径下的文件列表**
```
qshell [-d] dircache <DirCacheRootPath> <DirCacheResultFile>
```
比如，要获取`/Users/jemy/Temp4`目录下面的文件列表，则使用
```
qshell dircache /Users/jemy/Temp4 temp4.list.txt
```
其中`temp4.list.txt`是你保存列表结果的文件。列举的结果以如下格式组织：
```
文件相对于<DirCacheRootPath>的相对路径\t文件大小(单位字节)\t文件上次修改时间(单位100纳秒)
```
比如这样的：
```
rk_video_not_play.mp4	3985210	14206026340000000
rtl1.flv	10342916	14205959890000000
sync_demo/array_enumeration.png	5262899	13953255140000000
sync_demo/demo2.gif	2685960	13966636230000000
sync_demo/golang.png	149366	14010291080000000
```

**listbucket - 根据可选文件前缀来获取七牛空间中的文件列表**
```
qshell [-d] listbucket <Bucket> [<Prefix>] <ListBucketResultFile>
```
上面的`[Prefix]`表示这个`Prefix`参数是可选的，可以不设置，来获取空间中所有文件的列表。
如果设置了，则表示获取拥有指定文件前缀的所有文件列表。该列举结果的格式如下：
```
Key\tSize\tHash\tPutTime\tMimeType\tEndUser
```
比如：
```
hello.jpg	1710619	FlUqUK7zqbqm3NPwzq2q7TMZ-Ijs	14209629320769140	image/jpeg
hello.mp4	8495868	lns2dAHvO0qYseZFgDn3UqZlMOi-	14207312835630132	video/mp4
hhh	1492031	FjiRl_U0AeSsVCHXscCGObKyMy8f	14200176147531840	image/jpeg
jemygraw.jpg	1900176	FtmHAbztWfPEqPMv4t4vMNRYMETK	14208960018750329	application/octet-stream	QiniuAndroid
```
