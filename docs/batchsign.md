# 简介

`batchsign`命令用来根据资源的公开外链生成对应的私有外链，用于七牛私有空间的文件访问外链批量生成。

# 格式

```
qshell batchsign <UrlListFile> [<Deadline>]
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|可选参数|
|---------|-----------|---|
|UrlListFile|资源的公开访问外链列表，每行一个公开访问链接|N|
|Deadline|生成的私有访问外链的过期时间戳，如果不设置默认为一年有效期|Y|

# 示例

比如我们对文件`tosign.txt`里面的公开访问外链做签名。`tosign.txt`内容如下：

```
http://if-pri.qiniudn.com/camera.jpg
http://if-pri.qiniudn.com/camera.jpg?imageView2/0/w/100
```

使用

```
$ qshell batchsign tosign.txt
```

就能生成私有外链：

```
http://if-pri.qiniudn.com/camera.jpg?e=1473840685&token=TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:TnNXdt1Y4_jw-Xy0MF8vy9gF9dM=
http://if-pri.qiniudn.com/camera.jpg?imageView2/0/w/100&e=1473840685&token=TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:gjnUiiKUIOw7VQvJjYxXQLSybSM=
```

或者指定外链的有效期时间戳：

```
$ qshell batchsign tosign.txt 1473840685
```

这个时间戳可以用`d2ts`命令来生成。