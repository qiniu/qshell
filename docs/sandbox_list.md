# 简介
`sandbox list`（别名 `ls`）列出沙箱实例，支持按状态和元数据过滤。默认只列出运行中的沙箱。

# 格式
```
qshell sandbox list [-s <states>] [-m <metadata>] [-l <N>] [-f <pretty|json>]
qshell sbx ls [-s <states>] [-m <metadata>] [-l <N>] [-f <pretty|json>]
```

# 帮助文档
```
$ qshell sandbox list -h
$ qshell sandbox list --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `-s, --state`：按状态过滤，逗号分隔（可选值：running, paused）。默认为 running
- `-m, --metadata`：按元数据过滤（格式：key1=value1,key2=value2）
- `-l, --limit`：返回的最大数量
- `-f, --format`：输出格式，pretty（默认）或 json

# 输出列
| 列名 | 说明 |
|------|------|
| SANDBOX ID | 沙箱 ID |
| TEMPLATE ID | 模板 ID |
| ALIAS | 沙箱别名 |
| STARTED AT | 启动时间 |
| END AT | 结束时间 |
| STATE | 运行状态 |
| vCPUs | CPU 核数 |
| RAM MiB | 内存大小（MiB） |
| ENVD VERSION | envd 版本 |
| METADATA | 自定义元数据 |

# 示例
1. 列出运行中的沙箱（默认）
```
$ qshell sandbox list
$ qshell sbx ls
```

2. 列出所有状态的沙箱
```
$ qshell sandbox list -s running,paused
```

3. 按元数据过滤
```
$ qshell sandbox list -m user=alice,app=prod
```

4. 以 JSON 格式输出
```
$ qshell sandbox list -f json
```
