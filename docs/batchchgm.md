# 简介

`batchchgm`命令用来批量修改七牛空间中文件的MimeType。

# 格式

```
qshell batchchgm [-force] <Bucket> <KeyMimeMapFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|
|KeyMimeMapFile|文件名称和新的MimeType对的列表，每一行是`Key\tNewMimeType`格式，注意格式中间的Tab。|

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`-force`选项。

# 示例

比如我们要将空间`if-pbl`中的一些文件的MimeType修改为新的值。
那么提供的`KeyMimeMapFile`的内容有如下格式：

```
data/2015/02/01/bg.png	image/png
data/2015/02/01/pig.jpg	image/jpeg
```

在上面的列表中，`data/2015/02/01/bg.png`的新MimeType就是`image/png`，诸如此类。
把上面的内容保存在文件`tochange.txt`中，然后使用如下的命令：

```
$ qshell batchchgm if-pbl tochange.txt
```

如果执行过程中遇到任何错误，会输出到终端，如果没有的话，则没有任何输出。