# 简介
`buckets` 指令可以获取当前账号下所有的空间名称并输出。

# 格式
```
qshell buckets
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell buckets -h 

// 详细文档（此文档）
$ qshell buckets --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
无参数。

# 选项
- --region：指定需要列举所在区域的 bucket。
- --detail：打印 bucket 详情，如果无此选项则仅列举 bucket 名称，增加此选项后可依次展示 bucket 名、bucket 所在区域、bucket 文件数量、bucket 占用空间大小。

# 示例
简单使用
```
$ qshell buckets
```
输出：
```
bucket0
bucket1
bucket2
bucket3
bucket4
bucket7
```

列举所有区域 bucket 的详细信息
```
$ qshell buckets --detail
```
输出：
```
bucket0	z0	0	0(0B)
bucket1	z0	0	0(0B)
bucket2	z0	0	0(0B)
bucket3	z0	0	0(0B)
bucket4	z0	0	0(0B)
bucket7	z1	0	0(0B)
```

列举 z0 区域中 bucket 的详细信息
```
$ qshell buckets --region z0 --detail
```
输出：
```
bucket0	z0	0	0(0B)
bucket1	z0	0	0(0B)
bucket2	z0	0	0(0B)
bucket3	z0	0	0(0B)
bucket4	z0	0	0(0B)
```
