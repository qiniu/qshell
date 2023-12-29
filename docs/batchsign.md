# 简介
`batchsign` 命令用来根据资源的公开外链生成对应的私有外链，用于七牛私有空间的文件访问外链批量生成。

# 格式
```
qshell batchsign [-i <UrlListFile>] [-e <Deadline>]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell batchsign -h 

// 详细文档（此文档）
$ qshell batchsign --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- -i/--input-file：指定一个文件, 内容每行包含一个公开的外链。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：（【可选】）
```
<PublicUrl>   // 资源外链
```
- -o/--outfile：指定一个文件，把签名结果导入到此文件中【可选】
- -e/--deadline：接受一个过时的 deadline 参数，如果没有指定该参数，默认为 3600s 。【必选】 
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record；命令重新执行时，命令中所有任务会从头到尾重新执行；每个任务执行前会根据记录先查看当前任务是否已经执行，如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败则跳过不再重新执行。 【可选】

# 示例
比如我们对文件`tosign.txt`里面的公开访问外链做签名。`tosign.txt`内容如下：
```
http://if-pri.qiniudn.com/camera.jpg
http://if-pri.qiniudn.com/camera.jpg?imageView2/0/w/100
```
使用
```
$ qshell batchsign -i tosign.txt
```
就能生成私有外链：
```
http://if-pri.qiniudn.com/camera.jpg?e=1473840685&token=TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:TnNXdt1Y4_jw-Xy0MF8vy9gF9dM=
http://if-pri.qiniudn.com/camera.jpg?imageView2/0/w/100&e=1473840685&token=TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:gjnUiiKUIOw7VQvJjYxXQLSybSM=
```
或者指定外链的有效期时间戳：
```
$ qshell batchsign -i tosign.txt -e 1473840685
```
这个时间戳可以用`d2ts`命令来生成。

# 注意
如果没有指定输入文件，默认从标准输入读取内容
