# 简介
`cdnprefetch` 命令用来根据指定的文件访问列表来批量预取CDN的访问外链。

# 格式
```
qshell cdnprefetch [-i <UrlListFile>]
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name` 的情况下使用。

# 参数
无

# 选项
- -i/--input-file：接受一个文件参数，文件内容每行包含一个文件访问外链。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：【可选】
```
<Url> // <Url>：文件访问外链
```
- --qps：配置每秒预取的最大次数，默认不限制。【可选】
- -s/--size：每批预取的最大 Url 数，最大 50；默认 50。【可选】

# 示例
比如我们有如下内容的文件（`toprefetch.txt`），需要预取里面的外链
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
$ qshell cdnprefetch -i toprefetch.txt
```

就可以预取文件`toprefetch.txt`中的访问外链了。
