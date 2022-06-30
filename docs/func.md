# 简介
`func` 命令是封装 Go 语言的模板功能，你可已使用此命令进行模板的相关处理。

具体参考：
https://pkg.go.dev/text/template
https://masterminds.github.io/sprig

qshell 有些命令的回调使用模板方式实现，在使用之前可以使用此命令验证模板格式是否预期。

比如：在批量下载文件时，用户根据需要调整下载路径的回调；<ParamsJson> 是由工具内部传递包含有下载信息，用户仅需要编写 <FuncTemplate> 作为下载命令的参数，模板相当于回调函数，在函数中通过输出告诉 qshell 用户所期望的值(输出即 return 值)。

注：
输出的 [] 仅仅为了方便用户了解输出文字的边界，真正结果不包含 []。

# 格式
```
qshell func <ParamsJson> <FuncTemplate> [flags]
``` 

# 参数
- ParamsJson：Go 语言模板功能输入的参数 【必选】
- FuncTemplate：Go 语言的模板 【必选】

# 选项
无

# 示例
```
1. 截取 "this is a test" 中尾部的 test
$qshell func '{"name":"this is a test"}' '{{trimSuffix "test" .name}}'
[W]  output is insert [], and you should be careful with spaces etc.
[I]  [this is a ]

解析：
'{"name":"this is a test"}' 为参数，json 结构为字典
'{{trimSuffix "test" .name}}' 为模板，{{}}表示内部有函数调用，trimSuffix 为函数名， "test" 为需要去除的尾部字符串，. 表示当前参数(即上面字典)，name 表示字典中的 key，通过 .name 获取当前字典中键 name 对应的 key。
注意：' 和 { / } 之间没有其他符号，如果有则输出也会有(比如：'lala {{trimSuffix "test" .name}}'， 输出为：[lala this is a ])，其他符号包含空格等。


2. 截取 "this is a test" 中首部的 this
$qshell func '{"name":"this is a test"}' '{{trimPrefix "this" .name}}'
[W]  output is insert [], and you should be careful with spaces etc.
[I]  [ is a test]


3. 获取 "this is a test" 中从下标 0 到下标 5 的部分，包含下标 5。
qshell func '{"name":"this is a test"}' '{{substr 0 6 .name}}'
[W]  output is insert [], and you should be careful with spaces etc.
[I]  [this i]

```