# 简介
`mkbucket` 指令用来创建一个 bucket。

# 格式
```
qshell mkbucket <Bucket> [--region Region] [--private]
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数
- Bucket： 空间名称，可以为私有空间或者公开空间名称 【必选】

# 选项
- --region：指定创建 bucket 所在的区域；z0：华东，z1：华北，z2：华南，na0：北美，as0：东南亚(具体参考：https://developer.qiniu.com/kodo/1671/region-endpoint-fq)；默认为 z0。【可选】
- --private：是否创建私有空间；默认创建公有空间。【可选】

# 示例
在华北区域创建名为 my-bucket 的私有空间
```
$ qshell mkbucket my-bucket --region z1 --private
```