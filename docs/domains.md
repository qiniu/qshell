# 简介

`domains` 指令可以根据指定的空间参数获取和该空间关联的所有域名并输出。

# 格式

```
qshell domains <Bucket>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|-----|-----|
|Bucket|空间名称，可以为公开空间或者私有空间|

# 示例

获取空间`if-pbl`对应的所有的域名：

```
$ qshell domains if-pbl
```

输出：

```
if-pbl.qiniudn.com
7pn64c.com1.z0.glb.clouddn.com
7pn64c.com2.z0.glb.clouddn.com
7pn64c.com2.z0.glb.qiniucdn.com
```