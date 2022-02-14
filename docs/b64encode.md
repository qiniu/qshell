# 简介
该命令用来将一段字符串以`Base64编码`或`URL安全的Base64编码`格式进行编码。

# 格式
```
qshell b64encode [-s|--s] <DataToEncode>
```

# 参数
|    参数名    |            描述                | 可选参数 |
|--------------|--------------------------------|----------|
|      -s      |标志开启 urlsafe 的 base64 编码 | Y        |
| DataToDecode |待编码字符串                    | N        |

# 示例
```
$ qshell b64encode 'hello world'
aGVsbG8gd29ybGQ=
```
