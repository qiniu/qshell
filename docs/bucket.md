# 简介
`bucket` 指令用来获取 bucket 信息。

# 格式
```
qshell bucket <Bucket> 
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell bucket -h 

// 详细文档（此文档）
$ qshell bucket --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名称，可以为私有空间或者公开空间名称 【必选】

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
- Bucket：空间名。
- Region：区域信息；z0：华东，z1：华北，z2：华南，na0：北美，as0：东南亚(具体参考：https://developer.qiniu.com/kodo/1671/region-endpoint-fq)；默认为 z0。
- Private：是否为私有空间