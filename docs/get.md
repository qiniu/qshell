# 简介
`get` 用来从存储空间下载指定文件

# 格式
```
qshell get <Bucket> <Key> [-o <OutFile>]
``` 

# 参数
- Bucket：存储空间 【必选】
- Key：存储空间中的文件名字 【必选】

# 选项
- -o/--outfile：保存在本地的文件路径；不指定，保存在当前文件夹，文件名使用存储空间中的名字【可选】
- --domain：下载请求的 domain，qshell 下载使用 domain 优先级：1.domain(此选项) 2.bucket 配置域名(无需配置) 3.qshell 配置文件中 hosts 的 io(需要配置)，当优先级高的 domain 下载失败后会尝试使用优先级低的 domain 进行下载。【可选】
- --get-file-api: 当存储服务端支持 getfile 接口时才有效。【可选】
- --check-hash: 下载后检测本地文件和服务端文件 hash 的一致性。【可选】
- --enable-slice: 是否开启切片下载，需要注意 `--slice-file-size-threshold` 切片阈值选项的配置，只有开启切片下载，并且下载的文件大小大于切片阈值方会启动切片下载。默认不开启。【可选】
- --slice-size: 切片大小；当使用切片下载时，每个切片的大小；单位：B。默认为 4194304，也即 4MB。【可选】
- --slice-concurrent-count: 切片下载的并发度；默认为 10 【可选】
- --slice-file-size-threshold: 切片下载的文件阈值，当开启切片下载，并且文件大小大于此阈值时方会启用切片下载。【可选】
- --remove-temp-while-error: 当下载遇到错误时删除之前下载的部分文件缓存，默认为 `false` (不删除)【可选】

# 示例
1 把qiniutest空间下的文件test.txt下载到本地， 当前目录
```
$ qshell get qiniutest test.txt
```
下载完成后，在本地当前目录可以看到文件test.txt

2 把 `qiniutest` 空间下的文件 `test.txt` 下载到本地，以路径 `/Users/caijiaqiang/hah.txt` 保存。
```
$ qshell get qiniutest test.txt -o /Users/caijiaqiang/hah.txt
```
