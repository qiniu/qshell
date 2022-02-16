# 简介
`cdnrefresh` 命令用来根据指定的文件访问列表或者目录列表来批量更新CDN的缓存。

# 格式
刷新链接的命令格式：
```
qshell cdnrefresh [-i <UrlListFile>]
```

刷新目录的命令格式：
```
qshell cdnrefresh --dirs -i <DirListFile>
```

注意需要刷新的目录，必须以 `/` 结尾。如果没有制定输入文件 <UrlListFile> 默认从终端读取输入内容

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数和选项
### 刷新链接
##### i选项
接受一个文件参数，文件内容每行包含一个x需要进行刷新的文件访问外链。
文件每行具体格式如下：
```
<Url> // <Url>：文件访问外链
```
如果没有通过该选项指定该文件参数， 从标准输入读取内容。

### 刷新目录
接受一个文件参数，文件内容每行包含一个需要进行刷新的目录，目录必须以`/`结尾。
文件每行具体格式如下：
```
<FileUrl> // <FileUrl>：文件访问外链
```
如果没有通过该选项指定该文件参数， 从标准输入读取内容。

# 示例
比如我们有如下内容的文件（`torefresh.txt`），需要刷新里面的外链
```
http://if-pbl.qiniudn.com/hello1.txt
http://if-pbl.qiniudn.com/hello2.txt
http://if-pbl.qiniudn.com/hello3.txt
http://if-pbl.qiniudn.com/hello4.txt
http://if-pbl.qiniudn.com/hello5.txt
http://if-pbl.qiniudn.com/hello6.txt
http://if-pbl.qiniudn.com/hello7.txt
```

通过执行命令：
```
$ qshell cdnrefresh -i torefresh.txt
```

就可以刷新文件 `torefresh.txt` 中的访问外链了。
