# 简介

该命令用来将一段以`Base64编码`或`URL安全的Base64编码`编码的字符串解码。

# 格式

```
qshell b64decode [<UrlSafe>] <DataToDecode>
```

# 参数

|参数名|描述|可选参数|
|---------|-----------|----------|
|UrlSafe|指定是否以URL安全Base64编码格式解码，默认为true。|Y|
|DataToDecode|待解码字符串|N|

# 示例

我们可以解码七牛上传凭证的第三部分，即编码后的PutPolicy：

```
$ qshell b64decode 'eyJzY29wZSI6ImJiaW1nOnRlc3QucG5nIiwiZGVhZGxpbmUiOjE0MjcxODkxMzB9'
{"scope":"bbimg:test.png","deadline":1427189130}
```
