# 简介
`sandbox create`（别名 `cr`）创建一个新的沙箱实例并连接到其终端。当终端会话结束时，沙箱将被自动终止。

使用 `--detach` 模式时，创建沙箱后不连接终端，沙箱保持存活直到超时。

沙箱通过 keep-alive 机制保持存活，终端连接期间会自动续命，无需手动设置超时。

# 格式
```
qshell sandbox create <template> [-t <seconds>] [-d] [-m <metadata>]
qshell sbx cr <template> [-t <seconds>] [-d] [-m <metadata>]
```

# 帮助文档
```
$ qshell sandbox create -h
$ qshell sandbox create --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量，或在当前目录 `.env` 文件中设置。

# 参数
- `template`：模板 ID（必填）
- `-t, --timeout`：沙箱超时时间（秒）
- `-d, --detach`：创建沙箱但不连接终端，沙箱保持存活直到超时
- `-m, --metadata`：元数据键值对（格式：key1=value1,key2=value2）

# 示例
1. 创建沙箱
```
$ qshell sandbox create my-template
$ qshell sbx cr my-template
```

2. 设置超时时间
```
$ qshell sandbox create my-template --timeout 300
$ qshell sbx cr my-template -t 300
```

3. 分离模式创建（不连接终端，沙箱存活 5 分钟）
```
$ qshell sandbox create my-template -t 300 -d
$ qshell sbx cr my-template -t 300 -d
```

4. 添加元数据
```
$ qshell sandbox create my-template -m env=dev,team=backend
$ qshell sbx cr my-template -m env=dev,team=backend
```
