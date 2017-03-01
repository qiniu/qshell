# 简介

`reqid`命令用来解码一个七牛的自定义头部`X-Reqid`。

七牛会给每个请求添加一个自定义的头部，这个头部是根据一定的算法生成的，解码出来其实是一个日期，有时候我们需要这个日期来查询日志。

# 格式

```
qshell reqid <ReqIdToDecode>
```

# 参数

|参数名|描述|
|-----|-----|
|ReqIdToDecode|待解码的X-Reqid，注意是最后的一部分，比如对于`Reqid: shared.ffmpeg.62kAAIYB06brhtsT`，这里提供的值是`62kAAIYB06brhtsT`。|


# 示例

```
$ qshell reqid 62kAAIYB06brhtsT
2015-05-06/12-14
```