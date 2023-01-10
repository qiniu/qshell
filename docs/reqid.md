# 简介
`reqid` 命令用来解码一个七牛的自定义头部 `X-Reqid`。

七牛会给每个请求添加一个自定义的头部，这个头部是根据一定的算法生成的，解码出来其实是一个日期，有时候我们需要这个日期来查询日志。

# 格式
```
qshell reqid <ReqId>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell reqid -h 

// 详细文档（此文档）
$ qshell reqid --doc
```

# 鉴权
无

# 参数
- ReqId：待解码的 X-Reqid，注意是最后的一部分，比如对于 `Reqid: shared.ffmpeg.62kAAIYB06brhtsT`，这里提供的值是 `62kAAIYB06brhtsT`。 【必须】

# 示例
```
$ qshell reqid 62kAAIYB06brhtsT
2015-05-06/12-14
```