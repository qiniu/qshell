# 简介

`move`命令可以将一个空间中的文件移动到另外一个空间中，也可以对同一空间中的文件重命名。注意：移动文件仅支持在同一个帐号下面的空间中移动。
注意如果目标文件已存在空间中的时候，默认情况下，`move` 会失败，报错 `614 file exists`，如果一定要强制覆盖目标文件，可以使用选项 `-overwrite` 。

参考文档：[资源移动／重命名 (move)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/move.html)

# 格式

```
qshell move [-overwrite] <SrcBucket> <SrcKey> <DestBucket> <DestKey>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|--------|----------|
|SrcBucket|源空间名称|
|SrcKey|源文件名称|
|DestBucket|目标空间名称|
|DestKey|目标文件名称|

# 示例

1.将空间`if-pbl`中的`qiniu.jpg`移动到`if-pri`中

```
qshell move if-pbl qiniu.jpg if-pri qiniu.jpg
```

2.将空间`if-pbl`中的`qiniu.jpg`重命名为`2015/01/19/qiniu.jpg`

```
qshell move if-pbl qiniu.jpg if-pbl 2015/01/19/qiniu.jpg
```

3.将空间`if-pbl`中的`qiniu.jpg`移动到`if-pri`中，并命名为`2015/01/19/qiniu.jpg`

```
qshell move if-pbl qiniu.jpg if-pri 2015/01/19/qiniu.jpg
```

4.强制覆盖`if-pbl`中的已有文件`2015/01/19/qiniu.jpg`

```
qshell move -overwrite if-pbl qiniu.jpg if-pbl 2015/01/19/qiniu.jpg
```