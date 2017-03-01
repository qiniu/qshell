# 简介

该命令用来将一段字符串以`Base64编码`或`URL安全的Base64编码`格式进行编码。

# 格式

```
qshell b64encode [<UrlSafe>] <DataToEncode>
```

# 参数

|参数名|描述|可选参数|
|---------|-----------|----------|
|UrlSafe|指定是否以URL安全Base64编码格式编码，默认为true。|Y|
|DataToEncode|待编码字符串|N|

# 示例

```
$ qshell b64encode 'hello world'
aGVsbG8gd29ybGQ=
```
