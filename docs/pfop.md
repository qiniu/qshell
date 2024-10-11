# 简介
`pfop` 用来提交异步处理音视频请求， 比如视频转码，水印等， 打印服务端返回的 `PersistentID` 到标准输出上（终端）， 可以根据该 `PersistentID` 查询处理进度

参考文档：[pfop请求](http://developer.qiniu.com/code/v6/api/dora-api/pfop/pfop.html)

# 格式
```
qshell pfop [--pipeline <Pipeline>] <Bucket> <Key> <Fops>
qshell pfop [--pipeline <Pipeline>] --workflow-template-id <WorkflowTemplateID> <Bucket> <Key>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell pfop -h

// 详细文档（此文档）
$ qshell pfop --doc
```

# 参数
- Bucket：空间名，可以为公开空间或者私有空间【必选】
- Key：空间中文件的名称【必选】
- Fops：数据处理命令列表，以;分隔，可以指定多个数据处理命令。如果没有指定 `--workflow-template-id` 选项，则必选。
  如： `avthumb/mp4|saveas/cWJ1Y2tldDpxa2V5;avthumb/flv|saveas/cWJ1Y2tldDpxa2V5Mg==`，是将上传的视频文件同时转码成mp4格式和flv格式后另存。

# 选项
- -p/--pipeline：处理队列名称, 如果没有制定该选项，默认使用公有队列【可选】
- -u/--notify-url：处理结果通知接收 URL，七牛将会向你设置的 URL 发起 Content-Type: application/json 的 POST 请求。【可选】
- -y/--force：强制执行数据处理。当服务端发现 fops 指定的数据处理结果已经存在，那就认为已经处理成功，避免重复处理浪费资源。 增加此选项（--force），则可强制执行数据处理并覆盖原结果。【可选】
-    --workflow-template-id：工作流模版 ID【可选】
-    --type：任务类型【可选】

# 示例
1 把 qiniutest 空间下的文件 `test.avi` 转码成 `mp4` 文件，转码后的结果保存到 `qiniutest` 空间中
```
$ qshell pfop qiniutest test.avi 'avthumb/mp4'
```

返回 `PersistentID` 比如：
```
z1.5be96c32856db80b4be3d8b6
```

可以使用如下的命令查看处理进度
```
$ qshell prefop z1.5be96c32856db80b4be3d8b6
```
