#简介

`batchrename`命令用来为空间中的文件进行重命名。该操作发生在同一个空间中。

当然，如果这个时候源文件和目标文件名相同，那么复制会失败（这个操作其实没有意义），如果你实在想这样做，可以用`--overwrite`选项。


# 格式

```
qshell batchrename [--force] [--overwrite] [--success-list <SuccessFileName>] [--failure-list <failureFileName>] <Bucket> [-i <OldNewKeyMapFile>]
`****

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey`和`Name`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|

**i短选项**
接受一个文件参数, 内容为原文件名和目标文件名对的列表，注意这里目标文件名不可以和原文件名相同，否则对于这个文件来说就是重命名失败。每行的原文件名和目标文件名之间用`\t`分隔。如果没有指定该选项，默认从标准输入读取内容。

**success-list选项**
该选项指定一个文件，qshell会把操作成功的文件行导入到该文件

**failure-list选项**
该选项指定一个文件， qshell会把操作失败的文件行加上错误状态码，错误的原因导入该文件

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`--force`选项。

**overwrite选项**

默认情况下，如果批量重命名的文件列表中存在目标空间已有同名文件的情况，针对该文件的重命名会失败，如果希望能够强制覆盖目标文件，那么可以使用`--overwrite`选项。

# 示例

1.比如我们要将空间`if-pbl`里面的一些文件进行重命名，我们可以指定如下的`OldNewKeyMapFile`的内容：

```
2015/03/22/qiniu.png	test/qiniu.png
2015/photo.jpg	test/photo.jpg
```

上面，我们将`2015/03/22/qiniu.png`重命名为`test/qiniu.png`，诸如此类。
把这个内容保存到文件`torename.txt`中，然后使用如下的命令将所有的文件进行重命名。

```
$ qshell batchrename if-pbl -i torename.txt
```

2.如果不希望上面的重命名过程出现验证码提示，可以使用 `--force` 选项：

```
$ qshell batchrename --force if-pbl -i torename.txt
```

3. 对于重新命名的过程中，希望导入成功失败的文件，可以这样导出 

```
$ qshell batchrename if-pbl -i torename.txt --success-list success.txt --failure-list failure.txt
```

如果都重命名成功，success.txt的内容为：

```
2015/03/22/qiniu.png	test/qiniu.png
2015/photo.jpg	test/photo.jpg
```


# 注意 
如果没有指定输入文件的话， 会从标准输入读取内容
