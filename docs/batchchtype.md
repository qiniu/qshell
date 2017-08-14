# 简介

`batchchtype`命令用来为空间中的文件设置存储类型。该操作发生在同一个空间中。（将文件设置为**低频存储** 或者 **普通存储**，默认：文件为**普通存储**）


# 格式

```
qshell batchchtype [-force] <Bucket> <KeyFileTypeMapFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|
| KeyFileTypeMapFile |原文件名和存储类型的列表，存储类型用数字表示，1 为低频存储，0 为普通存储。每行的文件名和存储类型之间用`\t`分隔。|

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`-force`选项。

# 示例

1.比如我们要将空间`if-pbl`里面的一些文件改为低频存储，我们可以指定如下的`KeyFileTypeMapFile`的内容：

```
2015/03/22/qiniu.png	1
2015/photo.jpg	1
```

上面，我们将`2015/03/22/qiniu.png`文件设置为低频存储了，诸如此类。
把这个内容保存到文件`tochangetype.txt`中，然后使用如下的命令将 tochangetype.txt中所有的文件进行存储类型改变。

```
$ qshell batchchtype if-pbl tochangetype.txt
```

2.如果不希望上面的重命名过程出现验证码提示，可以使用 `-force` 选项：

```
$ qshell batchchtype -force if-pbl tochangetype.txt
```
