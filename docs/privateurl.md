# 简介

七牛空间分为公开空间和私有空间，无论是公开空间还是私有空间都对应一个默认的七牛的域名，这个域名也可以是用户自己的子域名。对于公开空间的资源访问，可以直接通过拼接域名和文件名的方式访问，而对私有空间中的资源，则还需要额外的授权操作。`privateurl`命令用来快速生成带签名的私有资源外链。 

# 格式

```
qshell privateurl <PublicUrl> [<Deadline>]
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数
|参数名|描述|可选参数|
|----------|-----------|----------|
|PublicUrl|资源的公开外链|N|
|Deadline|授权截至时间戳，单位秒|Y|

备注：

1. `Deadline`参数可以不指定，默认生成只有一个小时有效期的私有资源访问外链。
2. `Deadline`参数是一个单位为秒的Unix时间戳，可以使用`d2ts`命令生成。

# 示例

1.普通私有资源外链

```
$ qshell privateurl 'http://if-pri.qiniudn.com/beiyi.jpg'
```
结果：

```
http://if-pri.qiniudn.com/beiyi.jpg?e=1427613277&token=HCALkwxJcWd_8UlXCb6QWdA-pEZj1FXXSK0G1lMr:KrDZg1MGOmntVm5Hueny8l3JNjc=
```

2.带`fop`私有资源外链

```
$ qshell privateurl 'http://if-pri.qiniudn.com/beiyi.jpg?imageView2/0/w/600'

```
结果：

```
http://if-pri.qiniudn.com/beiyi.jpg?imageView2/0/w/600&e=1427613524&token=HCALkwxJcWd_8UlXCb6QWdA-pEZj1FXXSK0G1lMr:QzpohkbhnndlKFA2-YRGieVgGPE=
```
