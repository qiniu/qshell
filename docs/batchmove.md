# 简介

`batchmove`命令用来将一个空间中的文件批量移动到另一个空间，另外你可以在移动的过程中，给文件进行重命名。

当然，如果所指定的源空间和目标空间相同的话，如果这个时候源文件和目标文件名相同，那么移动会失败（这个操作其实没有意义）。
如果复制的目标空间中存在同名的文件，那么默认情况下针对该文件的移动操作也会失败，如果希望强制覆盖，可以指定`--overwrite`选项。

移动和复制的最大区别就是移动不保留原空间中的文件。

# 格式

```
qshell batchmove [-F <Delimiter>] [--force] [--overwrite] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>]<SrcBucket> <DestBucket> [-i <SrcDestKeyMapFile>]
```

# 帮助
```
qshell batchmove -h
```

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey`和`Name`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|SrcBucket|原空间名，可以为公开空间或私有空间|
|DestBucket|目标空间名，可以为公开空间或私有空间|

**i短选项**
接受一个文件参数，内容是原文件名和目标文件名对的列表，如果你希望目标文件名和原文件名相同的话，也可以不指定目标文件名，那么这一行就是只有原文件名即可。每行的原文件名和目标文件名之间用空白分隔（空格， \t, \n), 如果文件命中含有空格，可以使用-F选项指定别的分隔符。 如果没有指定该参数，默认从标准输入读取内容。

**success-list选项**
该选项指定一个文件，qshell会把操作成功的文件行导入到该文件

**failure-list选项**
该选项指定一个文件， qshell会把操作失败的文件行导入该文件

**force选项**

该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用`--force`选项。

**overwrite选项**

默认情况下，如果批量移动的文件列表中存在目标空间已有同名文件的情况，针对该文件的移动会失败，如果希望能够强制覆盖目标文件，那么可以使用`--overwrite`选项。

# 示例

1.我们将空间`if-pbl`中的一些文件移动到`if-pri`空间中去。如果是希望原文件名和目标文件名相同的话，可以这样指定`SrcDestKeyMapFile`的内容：

```
data/2015/02/01/bg.png
data/2015/02/01/pig.jpg
```

然后使用如下命令就可以以和原来相同的文件名从`if-pbl`移动到`if-pri`了。

```
$ qshell batchmove if-pbl if-pri -i tomove.txt
```


2.如果上面希望在移动的时候，对一些文件进行重命名，那么`SrcDestKeyMapFile`可以是这样：

```
data/2015/02/01/bg.png	background.png
data/2015/02/01/pig.jpg
```

从上面我们可以看到，你可以为你希望重命名的文件设置一个新的名字，不希望改变的就不用指定。

```
$ qshell batchmove if-pbl if-pri -i tomove.txt
```

3.如果不希望上面的移动过程出现验证码提示，可以使用 `--force` 选项：

```
$ qshell batchmove --force if-pbl if-pri -i tomove.txt
```

4.如果目标空间存在同名的文件，可以使用`--overwrite`选项来强制覆盖：

```
$ qshell batchmove --force --overwrite if-pbl if-pri -i tomove.txt
```

5. 假如文件名为tomove.txt, 文件包含了空格，所以必须使用其他的分隔符， 我们使用","

```
data/2015/02/01/bg.png,backg round.png
data/2015/02/01/pi g.jpg
```

```
$ qshell batchmove -i tomove.txt -F ',' if-pbl if-pri
```


# 注意
如果没有指定输入文件的话， 会从标准输入读取内容
