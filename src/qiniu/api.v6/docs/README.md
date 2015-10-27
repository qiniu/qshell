---
title: Go SDK 使用指南
---

此 Golang SDK 适用于所有 >=go1 版本，基于 [七牛云存储官方API](http://docs.qiniu.com) 构建。使用此 SDK 构建您的网络应用程序，能让您以非常便捷地方式将数据安全地存储到七牛云存储上。无论您的网络应用是一个网站程序，还是包括从云端（服务端程序）到终端（手持设备应用）的架构的服务或应用，通过七牛云存储及其 SDK，都能让您应用程序的终端用户高速上传和下载，同时也让您的服务端更加轻盈。

目录
----
- [概述](#overview)
- [安装](#install)
- [初始化](#setup)
	- [配置密钥](#setup-key)
- [上传文件](#io-put)
	- [上传流程](#io-put-flow)
	- [生成上传授权uptoken](#io-put-make-uptoken)
	- [上传代码](#io-put-upload-code)
	- [断点续上传、分块并行上传](#io-put-resumable)
	- [上传策略](#io-put-policy)
- [下载文件](#io-get)
	- [公有资源下载](#io-get-public)
	- [私有资源下载](#io-get-private)
	- [HTTPS支持](#io-get-https)
	- [断点续下载](#io-get-resumable)
- [资源操作](#rs)
	- [获取文件信息](#rs-stat)
	- [删除文件](#rs-delete)
	- [复制文件](#rs-copy)
	- [移动文件](#rs-move)
	- [批量操作](#rs-batch)
		- [批量获取文件信息](#rs-batch-stat)
		- [批量删除文件](#rs-batch-delete)
		- [批量复制文件](#rs-batch-copy)
		- [批量移动文件](#rs-batch-move)
		- [高级批量操作](#rs-batch-advanced)
- [数据处理接口](#fop-api)
	- [图像](#fop-image)
		- [查看图像属性](#fop-image-info)
		- [查看图片EXIF信息](#fop-exif)
		- [生成图片预览](#fop-image-view)
- [高级资源管理接口](#rsf-api)
	- [批量获得文件列表](#rsf-listPrefix)
- [贡献代码](#contribution)
- [许可证](#license)

----
<a name="overview"></a>

## 概述

七牛云存储的 GO 语言版本 SDK（本文以下称 GO-SDK）是对七牛云存储API协议的一层封装，以提供一套对于 GO 开发者而言简单易用的原生 GO 函数。GO 开发者在对接 GO-SDK 时无需理解七牛云存储 API 协议的细节，原则上也不需要对 HTTP 协议和原理做非常深入的了解，但如果拥有基础的 HTTP 知识，对于出错场景的处理可以更加高效。

GO-SDK 以开源方式提供。开发者可以随时从本文档提供的下载地址查看和下载 SDK 的源代码.

由于 GO 语言的通用性，GO-SDK 被设计为同时适合服务器端和客户端使用。服务端是指开发者自己的业务服务器，客户端是指开发者的客户终端。服务端因为有七牛颁发的 AccessKey/SecretKey，可以做很多客户端做不了的事情，比如删除文件、移动/复制文件等操作。一般而言，客服端操作文件需要获得服务端的授权。客户端上传文件需要获得服务端颁发的 [uptoken（上传授权凭证）](http://docs.qiniu.com/api/put.html#uploadToken)，客户端下载文件（包括下载处理过的文件，比如下载图片的缩略图）需要获得服务端颁发的 [dntoken（下载授权凭证）](http://docs.qiniu.com/api/get.html#download-token)。但开发者也可以将 bucket 设置为公开，此时文件有永久有效的访问地址，不需要业务服务器的授权，这对网站的静态文件（如图片、js、css、html）托管非常方便。

从 v5.0.0 版本开始，我们对 SDK 的内容进行了精简。所有管理操作，比如：创建/删除 bucket、为 bucket 绑定域名（publish）、设置数据处理的样式分隔符（fop seperator）、新增数据处理样式（fop style）等都去除了，统一建议到[开发者平台](https://portal.qiniu.com/)来完成。另外，此前服务端还有自己独有的上传 API，现在也推荐统一成基于客户端上传的工作方式。

从内容上来说，GO-SDK 主要包含如下几方面的内容：

* 公共库: api/conf
* 客户端上传文件：api/io
* 客户端断点续上传：api/resumable/io
* 数据处理：api/fop
* 服务端操作：api/auth/digest (授权), api/rs(资源操作, uptoken/dntoken颁发), api/rsf(批量获取文件列表)


<a name="install"></a>
## 1. 安装
在命令行下执行

	go get -u github.com/qiniu/api

<a name="setup"></a>
## 2. 初始化
<a name="setup-key"></a>
### 2.1 配置密钥

要接入七牛云存储，您需要拥有一对有效的 Access Key 和 Secret Key 用来进行签名认证。可以通过如下步骤获得：

1. [开通七牛开发者帐号](https://portal.qiniu.com/signup)
2. [登录七牛开发者自助平台，查看 Access Key 和 Secret Key](https://portal.qiniu.com/setting/key)

在获取到 Access Key 和 Secret Key 之后，您可以在您的程序中调用如下两行代码进行初始化对接, 要确保`ACCESS_KEY` 和 `SECRET_KEY` 在服务端调用 api/auth/digest,api/rs，api/rsf之前均已赋值：

```{go}
import . "github.com/qiniu/api/conf"

ACCESS_KEY = "<YOUR_APP_ACCESS_KEY>"
SECRET_KEY = "<YOUR_APP_SECRET_KEY>"
```

<a name="io-put"></a>
## 3. 上传文件

为了尽可能地改善终端用户的上传体验，七牛云存储首创了客户端直传功能。一般云存储的上传流程是：

    客户端（终端用户） => 业务服务器 => 云存储服务

这样多了一次上传的流程，和本地存储相比，会相对慢一些。但七牛引入了客户端直传，将整个上传过程调整为：

    客户端（终端用户） => 七牛 => 业务服务器

客户端（终端用户）直接上传到七牛的服务器，通过DNS智能解析，七牛会选择到离终端用户最近的ISP服务商节点，速度会比本地存储快很多。文件上传成功以后，七牛的服务器使用回调功能，只需要将非常少的数据（比如Key）传给应用服务器，应用服务器进行保存即可。

**注意**：如果您只是想要上传已存在您电脑本地或者是服务器上的文件到七牛云存储，可以直接使用七牛提供的 [qrsync](http://docs.qiniu.com/tools/qrsync.html) 上传工具。
文件上传有两种方式，一种是以普通方式直传文件，简称普通上传，另一种方式是断点续上传，断点续上传在网络条件很一般的情况下也能有出色的上传速度，而且对大文件的传输非常友好。



<a name="io-put-flow"></a>
### 3.1 上传流程

在七牛云存储中，整个上传流程大体分为这样几步：

1. 业务服务器颁发 [uptoken（上传授权凭证）](http://docs.qiniu.com/api/put.html#uploadToken)给客户端（终端用户）
2. 客户端凭借 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 上传文件到七牛
3. 在七牛获得完整数据后，发起一个 HTTP 请求回调到业务服务器
4. 业务服务器保存相关信息，并返回一些信息给七牛
5. 七牛原封不动地将这些信息转发给客户端（终端用户）

需要注意的是，回调到业务服务器的过程是可选的，它取决于业务服务器颁发的 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken)。如果没有回调，七牛会返回一些标准的信息（比如文件的 hash）给客户端。如果上传发生在业务服务器，以上流程可以自然简化为：

1. 业务服务器生成 uptoken（不设置回调，自己回调到自己这里没有意义）
2. 凭借 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 上传文件到七牛
3. 善后工作，比如保存相关的一些信息



<a name="io-put-make-uptoken"></a>
### 3.2 生成上传授权uptoken

uptoken是一个字符串,业务服务器根据(`rs.PutPolicy`)的结构体的各个参数来生成[uptoken](http://docs.qiniu.com/api/put.html#uploadToken)的代码如下:

调用如下代码前，请确保Access Key 和 Secret Key已经被正确初始化

```{go}
func uptoken(bucketName string) string {
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
		// CallbackUrl:  callbackUrl,
		// CallbackBody: callbackBody,
		// ReturnUrl:    returnUrl,
		// ReturnBody:   returnBody,
		// AsyncOps:     asyncOps,
		// EndUser:      endUser,
		// Expires:      expires,
	}
	return putPolicy.Token(nil)
}
```

参阅 `rs.PutPolicy` [policy参数](http://docs.qiniu.com/api/put.html#uploadToken-args)

<a name="io-put-upload-code"></a>
### 3.3 上传代码
上传文件到七牛（通常是客户端完成，但也可以发生在业务服务器）：

由于七牛的服务器支持自动生成key，所以本SDK提供的上传函数有两种展现方式，一种是有key的，一种是无key，让服务端自动生成key.

**注意： key必须采用utf8编码，如使用非utf8编码访问七牛云存储将反馈错误**

普通上传的文件和二进制，最后一个参数都是PutExtra类型，是用来细化上传功能用的，PutExtra的成员及其意义如下：

```{go}
type PutExtra struct {
	Params   map[string]string    //可选，用户自定义参数，必须以 "x:" 开头
	                              //若不以x:开头，则忽略
	MimeType string               //可选，当为 "" 时候，服务端自动判断
	Crc32    uint32
	CheckCrc uint32
	        // CheckCrc == 0: 表示不进行 crc32 校验
	        // CheckCrc == 1: 对于 Put 等同于 CheckCrc = 2；对于 PutFile 会自动计算 crc32 值
	        // CheckCrc == 2: 表示进行 crc32 校验，且 crc32 值就是上面的 Crc32 变量
}
```

直接上传内存中的数据, 代码:

```{go}
var err error
var ret io.PutRet
var extra = &io.PutExtra{
// Params:   params,
// MimeType: mieType,
// Crc32:    crc32,
// CheckCrc: CheckCrc,
}

// ret       变量用于存取返回的信息，详情见 io.PutRet
// uptoken   为业务服务器端生成的上传口令
// key       为文件存储的标识
// r         为io.Reader类型，用于从其读取数据
// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
err = io.Put(nil, &ret, uptoken, key, r, extra)

if err != nil {
	//上传产生错误
	log.Print("io.Put failed:", err)
	return
}

//上传成功，处理返回值
log.Print(ret.Hash, ret.Key)
```

参阅: `io.Put`, `io.PutExtra`

直接上传内存中的数据,且不提供key参数，此时key由七牛服务器自动生成, 代码:

```{go}
var err error
var ret io.PutRet
var extra = &io.PutExtra{
// Params:   params,
// MimeType: mieType,
// Crc32:    crc32,
// CheckCrc: CheckCrc,
}

// ret       变量用于存取返回的信息，详情见 io.PutRet
// uptoken   为业务服务器端生成的上传口令
// r         为io.Reader类型，用于从其读取数据
// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
err = io.PutWithoutKey(nil, &ret, uptoken, r, extra)

if err != nil {
	//上传产生错误
	log.Print("io.Put failed:", err)
	return
}

//上传成功，处理返回值
log.Print(ret.Hash, ret.Key)
```

参阅: `io.Put`, `io.PutExtra`

上传本地文件,代码:

```{go}
var err error
var ret io.PutRet
var extra = &io.PutExtra{
// Params:   params,
// MimeType: mieType,
// Crc32:    crc32,
// CheckCrc: CheckCrc,
}

// ret       变量用于存取返回的信息，详情见 io.PutRet
// uptoken   为业务服务器生成的上传口令
// key       为文件存储的标识
// localFile 为本地文件名
// extra     为上传文件的额外信息，详情见 io.PutExtra，可选
err = io.PutFile(nil, &ret, uptoken, key, localFile, extra)

if err != nil {
	//上传产生错误
	log.Print("io.PutFile failed:", err)
	return
}

//上传成功，处理返回值
log.Print(ret.Hash, ret.Key)
```

参阅: `io.PutFile`, `io.PutExtra`, `io.PutRet`

上传本地文件,且不提供key参数，此时key由七牛服务器自动生成代码:

```{go}
var err error
var ret io.PutRet
var extra = &io.PutExtra{
// Params:   params,
// MimeType: mieType,
// Crc32:    crc32,
// CheckCrc: CheckCrc,
}

// ret       变量用于存取返回的信息，详情见 io.PutRet
// uptoken   为业务服务器生成的上传口令
// localFile 为本地文件名
// extra     为上传文件的额外信息，详情见 io.PutExtra，可选
err = io.PutFileWithoutKey(nil, &ret, uptoken, localFile, extra)

if err != nil {
	//上传产生错误
	log.Print("io.PutFile failed:", err)
	return
}

//上传成功，处理返回值
log.Print(ret.Hash, ret.Key)
```

参阅: `io.PutFile`, `io.PutExtra`, `io.PutRet`

<a name="io-put-resumable"></a>

### 3.4 断点续上传、分块并行上传

除了基本的上传外，七牛还支持你将文件切成若干块（除最后一块外，每个块固定为4M大小），每个块可独立上传，互不干扰；每个分块块内则能够做到断点上续传。

断点续上传函数，最后一个选项是 resumable.io.PutExtra结构体，来细化上传用的，其成员及其含义如下：

```{go}
type PutExtra struct {
	Params     map[string]string // 可选。用户自定义参数，以"x:"开头 否则忽略
	MimeType   string            // 可选。
	ChunkSize  int               // 可选。每次上传的Chunk大小
	TryTimes   int               // 可选。尝试次数
	Progresses []BlkputRet       // 可选。上传进度
	Notify	   func(blkIdx int, blkSize int, ret *BlkputRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(blkIdx int, blkSize int, err error)
}
```

我们先看支持了断点上续传、分块并行上传的基本样例：
上传二进制流

```{go}
var err error
var ret rio.PutRet
var extra = &rio.PutExtra{
// Params:     params,
// MimeType:   mieType,
// ChunkSize:  chunkSize,
// TryTimes:   tryTimes,
// Progresses: progresses,
// Notify:     notify,
// NotifyErr:  NotifyErr,
}

// ret       变量用于存取返回的信息，详情见 resumable.io.PutRet
// uptoken   为业务服务器生成的上传口令
// key       为文件存储的标识
// r         为io.ReaderAt,用于读取数据
// fsize     数据总字节数
// extra     为上传文件的额外信息, 详情见 resumable.io.PutExtra
err = rio.Put(nil, &ret, uptoken, key, r, fsize, extra)

if err != nil {
	//上传产生错误
	log.Print("resumable.io.Put failed:", err)
	return
}

//上传成功，处理返回值
log.Print(ret.Hash)
```

参阅: `resumable.io.Put`, `resumable.io.PutExtra`, `rs.PutPolicy`

上传本地文件

```{go}
var err error
var ret rio.PutRet
var extra = &rio.PutExtra{
// Params:     params,
// MimeType:   mieType,
// ChunkSize:  chunkSize,
// TryTimes:   tryTimes,
// Progresses: progresses,
// Notify:     notify,
// NotifyErr:  NotifyErr,
}

// ret       变量用于存取返回的信息，详情见 resumable.io.PutRet
// uptoken   为业务服务器生成的上传口令
// key       为文件存储的标识
// localFile 为本地文件名
// extra     为上传文件的额外信息,可为空， 详情见 resumable.io.PutExtra
err = rio.PutFile(nil, &ret, uptoken, key, localFile, extra)

if err != nil {
	//上传产生错误
	log.Print("resumable.io.Put failed:", err)
	return
}

//上传成功，处理返回值
log.Print(ret.Hash)
```

参阅: `resumable.io.PutFile`, `resumable.io.PutExtra`, `rs.PutPolicy`

断点续上传的两个函数同样有两个wituoutKey函数`PutWithoutKey` `PutFileWithoutKey`与之对应，只是少了key参数，而这个key由我们的服务端自动生成。可通过函数返回的`PutRet`结构体获得key。
相比普通上传，断点上续传代码没有变复杂。基本上就只是将`io.PutExtra`改为`resumable.io.PutExtra`，`io.PutFile`改为`resumable.io.PutFile`。

但实际上 `resumable.io.PutExtra` 多了不少配置项，其中最重要的是两个回调函数：`Notify` 与 `NotifyErr`，它们用来通知使用者有更多的数据被传输成功，或者有些数据传输失败。在 `Notify` 回调函数中，比较常见的做法是将传输的状态进行持久化，以便于在软件退出后下次再进来还可以继续进行断点续上传。但不传入 `Notify` 回调函数并不表示不能断点续上传，只要程序没有退出，上传失败自动进行续传和重试操作。

<a name="io-put-policy"></a>
### 3.5 上传策略

[uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 实际上是用 AccessKey/SecretKey 进行数字签名的上传策略(`rs.PutPolicy`)，它控制则整个上传流程的行为。让我们快速过一遍你都能够决策啥：

* `Expires` 指定 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 有效期（默认1小时）。一个 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken) 可以被用于多次上传（只要它还没有过期）。
* `Scope` 限定客户端的权限。如果 `scope` 是 bucket，则客户端只能新增文件到指定的 bucket，不能修改文件。如果 `scope` 为 bucket:key，则客户端可以修改指定的文件。**注意： key必须采用utf8编码，如使用非utf8编码访问七牛云存储将反馈错误**
* `CallbackUrl` 设定业务服务器的回调地址，这样业务服务器才能感知到上传行为的发生。可选。
* `AsyncOps` 可指定上传完成后，需要自动执行哪些数据处理。这是因为有些数据处理操作（比如音视频转码）比较慢，如果不进行预转可能第一次访问的时候效果不理想，预转可以很大程度改善这一点。
* `ReturnBody` 可调整返回给客户端的数据包（默认情况下七牛返回文件内容的 `hash`，也就是下载该文件时的 `etag`）。这只在没有 `CallbackUrl` 时有效。
* `Escape` 为真（非0）时，表示客户端传入的 `CallbackParams` 中含有转义符。通过这个特性，可以很方便地把上传文件的某些元信息如 `fsize`（文件大小）、`ImageInfo.width/height`（图片宽度/高度）、`exif`（图片EXIF信息）等传给业务服务器。
* `DetectMime` 为真（非0）时，表示服务端忽略客户端传入的 `MimeType`，自己自行检测。

关于上传策略更完整的说明，请参考 [uptoken](http://docs.qiniu.com/api/put.html#uploadToken)。

<a name="io-get"></a>

## 4 下载文件

七牛云存储上的资源下载分为 公开资源下载 和 私有资源下载 。

私有（private）是 Bucket（空间）的一个属性，一个私有 Bucket 中的资源为私有资源，私有资源不可匿名下载。

新创建的空间（Bucket）缺省为私有，也可以将某个 Bucket 设为公有，公有 Bucket 中的资源为公有资源，公有资源可以匿名下载。

<a name="io-get-public"></a>

### 4.1 公有资源下载

如果在给bucket绑定了域名的话，可以通过以下地址访问。

	[GET] http://<domain>/<key>

其中<domain>是bucket所对应的域名。七牛云存储为每一个bucket提供一个默认域名。默认域名可以到[七牛云存储开发者平台](https://portal.qiniu.com/)中，空间设置的域名设置一节查询。用户也可以将自有的域名绑定到bucket上，通过自有域名访问七牛云存储。

**注意： key必须采用utf8编码，如使用非utf8编码访问七牛云存储将反馈错误**

<a name="io-get-private"></a>

### 4.2 私有资源下载

如果某个 bucket 是私有的，那么这个 bucket 中的所有文件只能通过一个的临时有效的 downloadUrl 访问：

	[GET] http://<domain>/<key>?token=<dnToken>

注意，尖括号不是必需，代表替换项。

其中 dntoken 是由业务服务器签发的一个[临时下载授权凭证](http://docs.qiniu.com/api/get.html#download-token)，deadline 是 dntoken 的有效期。dntoken不需要生成，GO-SDK 提供了生成完整 downloadUrl 的方法（包含了 dntoken），示例代码如下：

`downloadToken` 可以使用 SDK 提供的如下方法生成：

```{go}
func downloadUrl(domain, key string) string {
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}
```

生成 downloadUrl 后，服务端下发 downloadUrl 给客户端。客户端收到 downloadUrl 后，和公有资源类似，直接用任意的 HTTP 客户端就可以下载该资源了。唯一需要注意的是，在 downloadUrl 失效却还没有完成下载时，需要重新向服务器申请授权。

无论公有资源还是私有资源，下载过程中客户端并不需要七牛 GO-SDK 参与其中。

参阅: `rs.GetPolicy`, `rs.GetPolicy.MakeRequest`, `rs.MakeBaseUrl`

<a name="io-get-https"></a>

### 4.3 HTTPS支持

几乎所有七牛云存储 API 都同时支持 HTTP 和 HTTPS，但 HTTPS 下载有些需要注意的点。如果你的资源希望支持 HTTPS 下载，有如下限制：

1. 不能用 xxx.qiniudn.com 这样的二级域名，只能用 dn-xxx.qbox.me 域名。样例：https://dn-abc.qbox.me/1.txt
2. 使用自定义域名是付费的。我们并不建议使用自定义域名，但如确有需要，请联系我们的销售人员。

<a name="io-get-resumable"></a>
### 4.4 断点续下载

无论是公有资源还是私有资源，获得的下载 url 支持标准的 HTTP 断点续传协议。考虑到多数语言都有相应的断点续下载支持的成熟方法，七牛 GO-SDK 并不提供断点续下载相关代码。

<a name="rs"></a>
## 5. 资源操作

资源操作包括对存储在七牛云存储上的文件进行查看、复制、移动和删除处理。
该节调用的函数第一个参数都为 `logger`, 用于记录log, 如果无需求, 可以设置为nil. 具体接口可以查阅 `github.com/qiniu/rpc`

<a name="rs-stat"></a>
### 5.1 获取文件信息
函数`rs.Client.Stat`可获取文件信息。

```{go}
var ret rs.Entry
ret, err = rsCli.Stat(nil, bucket, key)
if err != nil {
	// 产生错误
	log.Println("rs.Stat failed:", err)
	return
}
// 处理返回值
log.Println(ret)
```

若有错误发生，则返回的err包含错误信息。若没错误返回的`ret`变量包含文件信息。
`ret`是为`rs.Entry`类型的结构体，其成员如下：

```{go}
type Entry struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Customer string `json:"customer"`
}
```

参阅: `rs.Entry`, `rs.Client.Stat`


<a name="rs-delete"></a>
### 5.2 删除文件
函数`rs.Client.Delete`可删除指定的文件。

```{go}
err = rsCli.Delete(nil, bucket, key)
if err != nil {
	// 产生错误
	log.Println("rs.Copy failed:", err)
	return
}
```

若无错误发生则返回的err为nil，否则err包含错误信息。
参阅: `rs.Client.Delete`

<a name="rs-copy"></a>
### 5.3 复制文件
函数`rs.Client.Copy`可用来复制文件。

```{go}
err = rsCli.Copy(nil, bucketSrc, keySrc, bucketDest, keyDest)
if err != nil {
	// 产生错误
	log.Println("rs.Copy failed:", err)
	return
}
```

参阅: `rs.Client.Move` `rs.Client.Copy`

<a name="rs-move"></a>
### 5.4 移动文件
函数`rs.Client.Move`可用来移动文件。

```{go}
err = rsCli.Move(nil, bucketSrc, keySrc, bucketDest, keyDest)
if err != nil {
	//产生错误
	log.Println("rs.Copy failed:", err)
	return
}
```

参阅: `rs.Client.Move`


<a name="rs-batch"></a>
### 5.5 批量操作
当您需要一次性进行多个操作时, 可以使用批量操作。

<a name="rs-batch-stat"></a>
#### 5.5.1 批量获取文件信息

函数`rs.Client.BatchStat`可批量获取文件信息。

```{go}
entryPathes := []rs.EntryPath{
	rs.EntryPath{
		Bucket: bucket1,
		Key:    key1,
	},
	rs.EntryPath{
		Bucket: bucket2,
		Key:    key2,
	},
}
var batchStatRets []rs.BatchStatItemRet
batchStatRets, err = rsCli.BatchStat(nil, entryPathes) // []rs.BatchStatItemRet, error
if err != nil {
	// 产生错误
	log.Println("rs.BatchStat failed:", err)
	return
}
// 处理返回值
for _, item := range batchStatRets {
	log.Println(item)
}
```

其中 `entryPathes`为`rs.EntryPath`结构体的数组切片。结构体`rs.EntryPath`中填写每个文件对应的bucket和key：

```{go}
type EntryPath struct {
	Bucket string
	Key string
}
```

`rs.BatchStat`会将文件信息(及成功/失败信息)，返回给由结构体`rs.BatchStatItemRet`组成的数组切片`batchStatRets`变量中。

```{go}
type BatchStatItemRet struct {
	Data  Entry  `json:"data"`
	Error string `json:"error"`
	Code  int    `json:"code"`
}
```

参阅: `rs.EntryPath`, `rs.BatchStatItemRet`, `rs.Client.BatchStat`

<a name="rs-batch-delete"></a>
#### 5.5.2 批量删除文件
函数`rs.Client.BatchDelete`可进行批量删除文件。

```{go}
entryPathes := []rs.EntryPath{
	rs.EntryPath{
		Bucket: bucket1,
		Key:    key1,
	},
	rs.EntryPath{
		Bucket: bucket2,
		Key:    key2,
	},
}
var batchDeleteRets []rs.BatchItemRet
batchDeleteRets, err = rsCli.BatchDelete(nil, entryPathes)
if err != nil {
	// 产生错误
	log.Println("rs.BatchDelete failed:", err)
	return
}
for _, item := range batchDeleteRets {
	// 遍历每个操作的返回结果
	log.Println(item.Code, item.Error)
}
```

和批量查看一样，`entryPathes`为`rs.EntryPath`结构体的数组切片。`rs.BatchDelete`会将删除操作的成功/失败信息返回给由结构体`rs.BatchItemRet`组成的数组切片`batchDeleteRets`变量中。其中`rs.BatchItemRet`结构体信息如下：

```{go}
type BatchItemRet struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
```

参阅: `rs.EntryPath`, `rs.Client.BatchDelete`, `rs.BatchItemRet`

<a name="rs-batch-copy"></a>
#### 5.5.3 批量复制文件
函数`rs.Client.BatchCopy`可进行批量复制文件。

```{go}
// 每个复制操作都含有源文件和目标文件
entryPairs := []rs.EntryPathPair{
	rs.EntryPathPair{
		Src: rs.EntryPath{
			Bucket: bucket1,
			Key:    key1,
		},
		Dest: rs.EntryPath{
			Bucket: bucket2,
			Key:    key2,
		},
	}, rs.EntryPathPair{
		Src: rs.EntryPath{
			Bucket: bucket3,
			Key:    key3,
		},
		Dest: rs.EntryPath{
			Bucket: bucket4,
			Key:    key4,
		},
	},
}
var batchCopyRets []rs.BatchItemRet
batchCopyRets, err = rsCli.BatchCopy(nil, entryPairs)
if err != nil {
	// 产生错误
	log.Println("rs.BatchCopy failed:", err)
	return
}
for _, item := range batchCopyRets {
	// 遍历每个操作的返回结果
	log.Println(item.Code, item.Error)
}
```

批量复制需要指明每个操作的源路径和目标路径，`entryPairs`是一个`rs.EntryPathPair`结构体的数组切片。结构体`rs.EntryPathPair`结构如下：

```{go}
type EntryPathPair struct {
	Src EntryPath
	Dest EntryPath
}
```

同样，`rs.BatchCopy`会将复制操作的成功/失败信息返回给由结构体`rs.BatchItemRet`组成的数组切片`batchCopyRets`变量中

参阅: `rs.BatchItemRet`, `rs.EntryPathPair`, `rs.Client.BatchCopy`

<a name="rs-batch-move"></a>
#### 5.5.4 批量移动文件
批量移动和批量很类似, 唯一的区别就是调用`rs.Client.BatchMove`

```{go}
// 每个复制操作都含有源文件和目标文件
entryPairs := []rs.EntryPathPair{
	rs.EntryPathPair{
		Src: rs.EntryPath{
			Bucket: bucket1,
			Key:    key1,
		},
		Dest: rs.EntryPath{
			Bucket: bucket2,
			Key:    key2,
		},
	}, rs.EntryPathPair{
		Src: rs.EntryPath{
			Bucket: bucket3,
			Key:    key3,
		},
		Dest: rs.EntryPath{
			Bucket: bucket4,
			Key:    key4,
		},
	},
}
var batchMoveRets []rs.BatchItemRet
batchMoveRets, err = rsCli.BatchMove(nil, entryPairs)
if err != nil {
	// 产生错误
	log.Println("rs.BatchMove failed:", err)
	return
}
for _, item := range batchMoveRets {
	// 遍历每个操作的返回结果
	log.Println(item.Code, item.Error)
}
```

参阅: `rs.EntryPathPair`, `rs.Client.BatchMove`

<a name="rs-batch-advanced"></a>
#### 5.5.5 高级批量操作
批量操作不仅仅支持同时进行多个相同类型的操作, 同时也支持不同的操作.

```{go}
ops := []string{
	rs.URIStat(bucket, key1),
	rs.URICopy(bucket, key1, bucket, key2), // 复制key1到key2
	rs.URIDelete(bucket, key1),             // 删除key1
	rs.URIMove(bucket, key2, bucket, key1), //将key2移动到key1
}

rets := new([]rs.BatchItemRet)
err = rsCli.Batch(nil, rets, ops)
if err != nil {
	// 产生错误
	log.Println("rs.Batch failed:", err)
	return
}
for _, ret := range *rets {
	log.Println(ret.Code, ret.Error)
}
```

参阅: `rs.URIStat`, `rs.URICopy`, `rs.URIMove`, `rs.URIDelete`, `rs.Client.Batch`

<a name="fop-api"></a>
## 6. 数据处理接口
七牛支持在云端对图像, 视频, 音频等富媒体进行个性化处理

<a name="fop-image"></a>
### 6.1 图像
<a name="fop-image-info"></a>
#### 6.1.1 查看图像属性
GO-SDK支持生成查看图片信息的URL，示意如下：

```{go}
func makeImageInfoUrl(imageUrl string) string {
	ii := fop.ImageInfo{}
	return ii.MakeRequest(imageUrl)
}
```

还可以已另一种方式，在程序中处理返回的图片信息：

```{go}
infoRet, err = ii.Call(nil, imageUrl)
if err != nil {
	// 产生错误
	log.Println("fop getImageInfo failed:", err)
	return
}
log.Println(infoRet.Height, infoRet.Width, infoRet.ColorModel,
	infoRet.Format)
```

参阅: `fop.ImageInfoRet`, `fop.ImageInfo`

<a name="fop-exif"></a>
#### 6.1.2 查看图片EXIF信息
同样，本SDK也支持直接生成查看exif的URL：

```{go}
func makeExifUrl(imageUrl string) string {
	e := fop.Exif{}
	return e.MakeRequest(imageUrl)
}
```

也可以在程序中处理exif的信息：

```{go}
exifRet, err = ie.Call(nil, imageUrl)
if err != nil {
	// 产生错误
	log.Println("fop getExif failed:", err)
	return
}

// 处理返回结果
for _, item := range exifRet {
	log.Println(item.Type, item.Val)
}
```

参阅: `fop.Exif`, `fop.ExifRet`, `fop.ExifValType`

<a name="fop-image-view"></a>
#### 6.1.3 生成图片预览
可以根据给定的文件URL和缩略图规格来生成缩略图的URL,代码：

```{go}
func makeViewUrl(imageUrl string) string {
	var view = fop.ImageView{
	// Mode    int    缩略模式
	// Width   int    Width = 0 表示不限定宽度
	// Height  int    Height = 0 表示不限定高度
	// Quality int    质量, 1-100
	// Format  string 输出格式，如jpg, gif, png, tif等等
	}
	return view.MakeRequest(imageUrl)
}
```

参阅: `fop.ImageView`

<a name="rsf-api"></a>
## 7. 高级资源管理接口(rsf)

<a name="rsf-listPrefix"></a>
### 7.1 批量获取文件列表
根据指定的前缀，获取对应前缀的文件列表,正常使用情景如下：

```{go}
func listAll(l rpc.Logger, rs *rsf.Client, bucketName string, prefix string) {

	var entries []rsf.ListItem
	var marker = ""
	var err error
	var limit = 1000

	for err == nil {
		entries, marker, err = rs.ListPrefix(l, bucketName,
			prefix, marker, limit)
		for _, item := range entries {
			// 处理 item
			log.Print("item:", item)
		}
	}
	if err != io.EOF {
		// 非预期的错误
		log.Print("listAll failed:", err)
	}
}
```

参阅: `rsf.ListPreFix`


<a name="contribution"></a>
## 8. 贡献代码

1. Fork
2. 创建您的特性分支 (`git checkout -b my-new-feature`)
3. 提交您的改动 (`git commit -am 'Added some feature'`)
4. 将您的修改记录提交到远程 `git` 仓库 (`git push origin my-new-feature`)
5. 然后到 github 网站的该 `git` 远程仓库的 `my-new-feature` 分支下发起 Pull Request

<a name="license"></a>
## 9. 许可证

Copyright (c) 2013 qiniu.com

基于 MIT 协议发布:

* [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)

