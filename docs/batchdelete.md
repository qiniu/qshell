# 简介

`batchdelete`命令用来根据一个七牛空间中的文件名列表来批量删除空间中的这些文件。

# 格式

```
qshell batchdelete [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] <Bucket> [-i <KeyListFile>]
```

# 帮助
```
qshell batchdelete -h
```

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey`和`Name`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|
|KeyListFile|文件列表文件，该列表文件只要保证第一列是文件名即可，每个列用`\t`分隔，可以直接使用`listbucket`的结果。|

**success-list选项**
该选项指定一个文件，qshell会把操作成功的文件行导入到该文件

**failure-list选项**
该选项指定一个文件， qshell会把操作失败的文件行加上错误状态码，错误的原因导入该文件

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`--force`选项。

# 示例

1.指定要删除的文件列表`todelete.txt`进行删除：

```
a.jpg
test/b.jpg
```

```
$ qshell batchdelete if-pbl -i todelete.txt
```

2.删除空间`if-pbl`中的所有文件：

```
$ qshell listbucket if-pbl -o if-pbl.list.txt
$ qshell batchdelete --force if-pbl -i if-pbl.list.txt
```

3. 如果希望导出成功和失败的文件列表

```
$ qshell batchdelete if-pbl -i if-pbl.list.txt --success-list success.txt --failure-list failure.txt
```
