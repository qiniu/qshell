# 简介

`m3u8replace`命令用来修改或删除七牛空间中m3u8播放列表中引用的切片路径中的域名。

# 格式

```
qshell m3u8replace <Bucket> <M3u8Key> [<NewDomain>]
``` 

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|可选参数|
|--------|--------|-------|
|Bucket|m3u8文件所在空间，可以为公开空间或私有空间|N|
|M3u8Key|m3u8文件的名字|N|
|NewDomain|引用切片的域名，如果不指定的话，则m3u8文件中引用切片使用相对路径，效果等同于转码时指定`noDomain/1`|Y|

# 示例

1.清除m3u8播放列表中切片引用路径中的域名，等同于转码时指定`noDomain/1`

```
qshell m3u8replace if-pbl qiniu.m3u8
```

2.替换m3u8播放列表中切片引用路径中的域名，把旧的换成新的。

```
qshell m3u8replace if-pbl qiniu.m3u8 http://hls.example.com
```
 