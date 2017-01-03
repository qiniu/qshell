# 简介

`batchcopy`命令用来将一个空间中的文件批量复制到另一个空间，另外你可以在复制的过程中，给文件进行重命名。

# 格式

```
qshell batchcopy <SrcBucket> <DestBucket> <SrcDestKeyMapFile>
```
# 参数

|参数名|描述|
|---------|-----------|
|SrcBucket|原空间名，可以为公开空间或私有空间|
|DestBucket|目标空间名，可以为公开空间或私有空间|
|SrcDestKeyMapFile|原文件名和目标文件名对的列表，如果你希望目标文件名和原文件名相同的话，也可以不指定目标文件名，那么这一行就是只有原文件名即可。每行的原文件名和目标文件名之间用`\t`分隔。|

# 示例

1.我们将空间`if-pbl`中的一些文件复制到`if-pri`空间中去。如果是希望原文件名和目标文件名相同的话，可以这样指定`SrcDestKeyMapFile`的内容：
```
data/2015/02/01/bg.png
data/2015/02/01/pig.jpg
```
然后使用如下命令：
```
$ qshell batchcopy if-pbl if-pri tocopy.txt
```
那么上面的文件就以和原来相同的文件名从`if-pbl`复制到`if-pri`了。

2.如果上面希望在复制的时候，对一些文件进行重命名，那么`SrcDestKeyMapFile`可以是这样：
```
data/2015/02/01/bg.png	background.png
data/2015/02/01/pig.jpg
```
从上面我们可以看到，你可以为你希望重命名的文件设置一个新的名字，不希望改变的就不用指定。
然后使用命令：
```
$ qshell batchcopy if-pbl if-pri tomove.txt
```
就可以将文件复制过去了。