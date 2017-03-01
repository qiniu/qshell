# 简介

`cdnprefetch`命令用来根据指定的文件访问列表来批量预取CDN的访问外链。

# 格式

```
qshell cdnprefetch <UrlListFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|UrlListFile|需要进行预取的文件访问外链列表，每行一个访问外链|

# 示例

比如我们有如下内容的文件，需要预取里面的外链

```
http://if-pbl.qiniudn.com/hello1.txt
http://if-pbl.qiniudn.com/hello2.txt
http://if-pbl.qiniudn.com/hello3.txt
http://if-pbl.qiniudn.com/hello4.txt
http://if-pbl.qiniudn.com/hello5.txt
http://if-pbl.qiniudn.com/hello6.txt
http://if-pbl.qiniudn.com/hello7.txt
```

```
$ qshell cdnprefetch toprefetch.txt
```

就可以预取文件`toprefetch.txt`中的访问外链了。