#简介

`batchrename`命令用来为空间中的文件进行重命名。该操作发生在同一个空间中。

当然，如果这个时候源文件和目标文件名相同，那么复制会失败（这个操作其实没有意义），如果你实在想这样做，可以用`-overwrite`选项。


# 格式

```
qshell batchrename [-force] [-overwrite] <Bucket> <OldNewKeyMapFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|
|OldNewKeyMapFile|原文件名和目标文件名对的列表，注意这里目标文件名不可以和原文件名相同，否则对于这个文件来说就是重命名失败。每行的原文件名和目标文件名之间用`\t`分隔。|

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`-force`选项。

**overwrite选项**

默认情况下，如果批量重命名的文件列表中存在目标空间已有同名文件的情况，针对该文件的重命名会失败，如果希望能够强制覆盖目标文件，那么可以使用`-overwrite`选项。

# 示例

1.比如我们要将空间`if-pbl`里面的一些文件进行重命名，我们可以指定如下的`OldNewKeyMapFile`的内容：

```
2015/03/22/qiniu.png	test/qiniu.png
2015/photo.jpg	test/photo.jpg
```

上面，我们将`2015/03/22/qiniu.png`重命名为`test/qiniu.png`，诸如此类。
把这个内容保存到文件`torename.txt`中，然后使用如下的命令将所有的文件进行重命名。

```
$ qshell batchrename if-pbl torename.txt
```

2.如果不希望上面的重命名过程出现验证码提示，可以使用 `-force` 选项：

```
$ qshell batchrename -force if-pbl torename.txt
```