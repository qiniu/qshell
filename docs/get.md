# 简介
`get` 用来从存储空间下载指定的文件。

# 格式
```
qshell get <Bucket> <Key> [-o <OutFile>]
``` 

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell get -h 

// 详细文档（此文档）
$ qshell get --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名。 【必选】
- Key：存储空间中的文件名字。 【必选】

## 注：
1. 使用 bucket 绑定的源站域名和七牛源站域名下载资源，这部分下载产生的流量会生成存储源站下载流量的计费，请注意，这部分计费不在七牛 CDN 免费 10G 流量覆盖范围，具体域名使用参考：--domain 选项。

# 选项
- -o/--outfile：保存在本地的文件路径；不指定，保存在当前文件夹，文件名使用存储空间中的名字【可选】
- --domain：指定下载请求的域名，当指定了下载域名则仅使用此下载域名进行下载；默认为空，此时 qshell 下载使用域名的优先级：1.bucket 绑定的 CDN 域名(qshell 内部查询，无需配置) 2.bucket 绑定的源站域名(qshell 内部查询，无需配置) 3. 七牛源站域名(qshell 内部查询，无需配置)，当优先级高的域名下载失败后会尝试使用优先级低的域名进行下载。【可选】
- --get-file-api: 当存储服务端支持 getfile 接口时才有效。【可选】
- --public：空间是否为公开空间；为 `true` 时为公有空间，公有空间下载时不会对下载 URL 进行签名，可以提升 CDN 域名性能，默认为 `false`（私有空间），已废弃【可选】
- --check-size: 下载后检测本地文件和服务端文件 size 的一致性。【可选】
- --check-hash: 下载后检测本地文件和服务端文件 hash 的一致性。【可选】
- --enable-slice: 是否开启切片下载，需要注意 `--slice-file-size-threshold` 切片阈值选项的配置，只有开启切片下载，并且下载的文件大小大于切片阈值方会启动切片下载。默认不开启。【可选】
- --slice-size: 切片大小；当使用切片下载时，每个切片的大小；单位：B。默认为 4194304，也即 4MB。【可选】
- --slice-concurrent-count: 切片下载的并发度；默认为 10 【可选】
- --slice-file-size-threshold: 切片下载的文件阈值，当开启切片下载，并且文件大小大于此阈值时方会启用切片下载。【可选】
- --remove-temp-while-error: 当下载遇到错误时删除之前下载的部分文件缓存，默认为 `false` (不删除)【可选】

注：
如果使用的是 CDN 域名，且 CDN 域名开启了图片优化中的图片自动瘦身功能时，下载文件的信息和七牛服务端记录的文件信息不一致，此时下载不要使用 --check-size 和 --check-hash 选项，否则下载会失败。

# 示例
1 把 `qiniutest` 空间下的文件 test.txt 下载到本地的当前目录：
```
$ qshell get qiniutest test.txt
```
下载完成后，在本地当前目录可以看到文件 test.txt

2 把 `qiniutest` 空间下的文件 `test.txt` 下载到本地，以路径 `/Users/caijiaqiang/hah.txt` 进行保存。
```
$ qshell get qiniutest test.txt -o /Users/caijiaqiang/hah.txt
```
