# 简介
`get` 用来从存储空间下载文件，该空间不需要绑定域名也可以下载

# 格式
```
qshell get <Bucket> <Key> [-o <OutFile>]
``` 

# 参数
- Bucket：存储空间 【必选】
- Key：存储空间中的文件名字 【必选】

# 选项
- -o：保存在本地的文件路径；不指定，保存在当前文件夹，文件名使用存储空间中的名字【可选】
- --domain：下载请求的 domain 信息。【可选】
- --get-file-api: 公有云无效，当私有云支持 getfile 接口时有效【可选】
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
