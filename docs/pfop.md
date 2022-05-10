# 简介

`pfop` 用来提交异步处理音视频请求， 比如视频转码，水印等， 打印服务端返回的persistentID到标准输出上（终端）， 可以根据该persistentId查询处理进度

# 格式

```
qshell pfop [--pipeline <Pipeline>] <Bucket> <Key> <fopCommand>
``` 

# 参数

|参数名|描述|
|--------|--------|
|pipeline|处理队列名称, 如果没有制定该选项，默认使用公有队列 |

# 示例

1. 把qiniutest空间下的文件test.avi文件转码成mp4文件, 转码后的结果保存到qiniutest空间中

```
$ qshell pfop qiniutest test.avi 'avthumb/mp4'
```

返回persistentId比如：

```
z1.5be96c32856db80b4be3d8b6
```

可以使用如下的命令查看处理进度

```
$ qshell prefop z1.5be96c32856db80b4be3d8b6
```
