# 简介

`get` 用来从存储空间下载文件，该空间不需要绑定域名也可以下载

# 格式

```
qshell get <Bucket> <Key> [-o <OutFile>]
``` 

# 参数

|参数名|描述|
|--------|--------|
|Bucket |存储空间|
|Key|存储空间中的文件名字|
|OutFile|保存在本地的名字，不指定，默认使用存储空间中的名字|

# 示例

1. 把qiniutest空间下的文件test.txt下载到本地， 当前目录

```
$ qshell get qiniutest test.txt
```
下载完成后，在本地当前目录可以看到文件test.txt

2. 把qiniutest空间下的文件test.txt下载到本地，以路径/Users/caijiaqiang/hah.txt保存

```
$ qshell get qiniutest test.txt -o /Users/caijiaqiang/hah.txt
```
