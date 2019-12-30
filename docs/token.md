# 简介

`token`是用来计算Qbox token, Qiniu Token, Upload Token的命令。

# 格式

```
qshell token qbox [--access-key <AccessKey> --secret-key <SecretKey>] [--http-body <Body>] [--content-type <Type>] <Url>
qshell token qiniu [--access-key <AccessKey> --secret-key <SecretKey>] [--http-body <Body>] [--content-type <Type>] [--method <Method>] <Url>
qshell token upload <PutPolicyConfigFile>
```

# 配置

计算upload token需要上传策略配置文件， 配置文件支持的全部参数如下：

```
{
	scope                     : "",
	deadline                  : <Unix时间戳>,
	isPrefixalScope         : [1|0],
	insertOnly               : [1|0],
	detectMime               : [1|0],
	fsizeLimit               : <限制文件的大小>,
	mimeLimit                : <限制允许上传的文件MimeType>,
	saveKey                  : <自定义上传的资源名>,
	callbackFetchKey        : "",
	callbackUrl              : <回调地址>,
	callbackHost             : <回调Host>,
	callbackBody             : <回调Body>,
	callbackBodyType        : <回调ContentType>,
	returnUrl                : "",
	returnBody               : "",
	persistentOps            : "",
	persistentNotifyUrl     : "",
	persistentPipeline       : "",
	endUser                  : "",
	deleteAfterDays         : "",
	fileType                 : [1|0]
}
```

上传策略的各个参数[详解](https://developer.qiniu.com/kodo/manual/1206/put-policy) 。


# 例子

一.七牛的stat接口，可以获取一个bucket中的文件信息， 详细[文档](https://developer.qiniu.com/kodo/api/1308/stat)。
这个接口需要计算QBox Token, 那么我们可以使用qshell token qbox <Url>这个命令来计算。
Url的格式为： `http://<Host>/<Path>`

从接口文档，我们可以看到Host是"rs.qiniu.com",  Path是"/stat/<EncodedEntryURI>", <EncodeEntryURI>的计算方式参考
[文档](https://developer.qiniu.com/kodo/api/1276/data-format), 可以通过qshell b64encode来计算这个值。

假设我们需要拿到文件`test.mov`的信息， 这个文件存储在七牛的存储空间`tonypublic`中，那么，
第一步， 计算<EncodeEntryURI>：

```
$ qshell b64encode "tonypublic:test.mov"
```

得到`dG9ueXB1YmxpYzp0ZXN0Lm1vdg==`

第二步， 计算<Url>, URL = `http://rs.qiniu.com/stat/dG9ueXB1YmxpYzp0ZXN0Lm1vdg==`

第三步， 计算Qbox Token:

```
$ qshell token qbox http://rs.qiniu.com/stat/dG9ueXB1YmxpYzp0ZXN0Lm1vdg==
```

得到qbox token: `QBox 3-pH6WfqAXTwzgG2s3FNMUW0NtkUu5cJLQCfU3Hd:d5rmqx_xsKILoNa61qDHMxUMtp8=`

第四步， 通过curl 发送http请求，拿到test.mov的信息:

```
$ url 'http://rs.qiniu.com/stat/dG9ueXB1YmxpYzp0ZXN0Lm1vdg==' -H "Authorization: QBox 3-pH6WfqAXTwzgG2s3FNMUW0NtkUu5cJLQCfU3Hd:d5rmqx_xsKILoNa61qDHMxUMtp8="
```
得到该文件的信息如下:

```json
{"fsize":94633760,"hash":"lhuYnUQEvCavdrNrrz82nEWSSqsB","md5":"LPkxXhSLb5fb9fxrLUghkA==","mimeType":"text/html","putTime":15289618585925391,"type":0}
```




二.计算上传token

上传策略配置文件upload.conf内容为：

```json
{
    "scope": "tonypublic"
}
```

可以使用如下的命令计算：

```
$ qshell token upload upload.conf
```

得到`UpToken 3-pH6WfqAXTwzgG2s3FNMUW0NtkUu5cJLQCfU3Hd:EYUNznmCcnlhFU5a126AKwmoHgE=:eyJzY29wZSI6InRvbnlwdWJsaWMiLCJkZWFkbGluZSI6MTU0NDQzMjY5MH0=`
