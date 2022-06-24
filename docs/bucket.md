# 简介
`bucket` 指令用来获取 bucket 信息。

# 格式
```
qshell bucket <Bucket> 
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数
- Bucket： 空间名称，可以为私有空间或者公开空间名称 【必选】

# 选项
无

# 示例
获取 my-bucket 空间的信息
```
$ qshell bucket my-bucket
```

输出：
```
Bucket    :my-bucket
Region    :z0
Private   :true
```

输出字段说明：
- Bucket：bucket 的名称。
- Region：区域信息；z0：华东，z1：华北，z2：华南，na0：北美，as0：东南亚(具体参考：https://developer.qiniu.com/kodo/1671/region-endpoint-fq)；默认为 z0。
- Private：是否为私有空间