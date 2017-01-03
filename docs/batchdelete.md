# 简介

`batchdelete`命令用来根据一个七牛空间中的文件名列表来批量删除空间中的这些文件。

# 格式

```
qshell batchdelete [-force] <Bucket> <KeyListFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|
|KeyListFile|文件列表文件，该列表文件只要保证第一列是文件名即可，每个列用`\t`分隔，可以直接使用`listbucket`的结果。|

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`-force`选项。

# 示例

1.指定要删除的文件列表`todelete.txt`进行删除：

```
a.jpg
test/b.jpg
```

```
$ qshell batchdelete if-pbl todelete.txt
```

2.删除空间`if-pbl`中的所有文件：

```
$ qshell listbucket if-pbl if-pbl.list.txt
$ qshell batchdelete -force if-pbl if-pbl.list.txt
```
