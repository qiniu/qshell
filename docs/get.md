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
- --domain：下载请求的 domain 信息。【可选】
- --get-file-api: 当存储服务端支持 getfile 接口时才有效。【可选】
- --check-hash: 下载后检测本地文件和服务端文件 hash 的一致性。【可选】
- --big-file-enable-slice: 当文件大于 40M 时采用分片下载，每个分片大小为 4M ，10 并发下载。
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
