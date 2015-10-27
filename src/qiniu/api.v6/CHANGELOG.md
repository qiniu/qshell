## CHANGE LOG

### v6.0.6

2014-6-30 issues [#114](https://github.com/qiniu/api/pull/114)

- [#111] 修改mimetype
- [#112] 去除chunked 上传
- [#115] 统一user agent


### v6.0.5

2014-6-06 issues [#108](https://github.com/qiniu/api/pull/108)

- [#107] 增加pipeline, gofmt


### v6.0.4

2014-1-19 issues [#103](https://github.com/qiniu/api/pull/103)

- 文档及一些细节调整


### v6.0.3

2013-10-08 issues [#98](https://github.com/qiniu/api/pull/98)

- 更新断点续上传到新接口
- 修订文档
- 更新断点续上传测试用例

### v6.0.2

2013-08-03 issue [#88](https://github.com/qiniu/api/pull/88)

- bugfix: GetPolicy.MakeRequest
- 增加私有资源的图片缩略图样例


### v6.0.1

2013-07-01 issue [#77](https://github.com/qiniu/api/pull/77)

- 整理 sdk docs


### v6.0.0

2013-06-29 issue [#73](https://github.com/qiniu/api/pull/73)

- 遵循 [sdkspec v6.0.2](https://github.com/qiniu/sdkspec/tree/v6.0.2)
  - UserAgent 支持。增加 conf.SetUserAgent
  - io.Put/PutFile 调整为基于 up.qiniu.com 的协议，extra *PutExtra 参数可以为 nil
  - io.Put/PutFile 支持支持 key = UNDEFINED_KEY，这样服务端将自动生成 key 并返回
  - io.Put/PutFile 支持自定义的 "x:" 参数(io.PutExtra.Params)、支持 Crc 检查
  - 增加 rsf api 的支持
- 修复 PutPolicy.Token 调用多次出错的 bug


### v5.0.0

2013-06-11 issue [#62](https://github.com/qiniu/api/pull/62)

- 遵循 [sdkspec v1.0.2](https://github.com/qiniu/sdkspec/tree/v1.0.2)
  - rs.GetPolicy 删除 Scope，也就是不再支持批量下载的授权。
  - rs.New, PutPolicy.Token, GetPolicy.MakeRequest 增加 mac *digest.Mac 参数。
- 初步整理了 sdk 使用文档。


### v0.9.1

2013-05-28 issue [#56](https://github.com/qiniu/api/pull/56)

- 修复 go get github.com/qiniu/api 失败的错误
- 遵循 [sdkspec v1.0.1](https://github.com/qiniu/sdkspec/tree/v1.0.1)
  - io.GetUrl 改为 rs.MakeBaseUrl 和 rs.GetPolicy.MakeRequest
  - rs.PutPolicy 支持 ReturnUrl, ReturnBody, CallbackBody；将 Customer 改为 EndUser
- 增加 github.com/api/url: Escape/Unescape


### v0.9.0

2013-04-08 issue [#33](https://github.com/qiniu/api/pull/33)

- 更新API文档
- 增加断点续上传resumable/io功能
- 移除bucket相关的增加/删除/列出所有bucket等管理操作，推荐到七牛云存储开发者后台中使用这些功能。
