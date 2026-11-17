# 简介
`sandbox list`（别名 `ls`）列出沙箱实例，支持按状态和元数据过滤。默认只列出运行中的沙箱。

# 格式
```
qshell sandbox list [-s <states>] [-m <metadata>] [-l <N>] [--all | --next-token <cursor>] [-f <pretty|json>]
qshell sbx ls [-s <states>] [-m <metadata>] [-l <N>] [--all | --next-token <cursor>] [-f <pretty|json>]
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
- `-l, --limit`：每页的最大数量，省略时使用服务端默认值；使用 `--all` 时不限制总数量
- `--all`：读取全部匹配页后统一输出
- `--next-token`：从上一页的游标继续，不能与 `--all` 同用
- `-f, --format`：输出格式，pretty（默认）或 json

默认读取一页；还有下一页时，游标以 `Next token: ...` 输出到 stderr。JSON 的 stdout 保持沙箱数组，便于继续交给现有脚本处理。任一页查询失败时返回非零退出码；`--all` 不会把部分结果作为完整查询输出。

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

5. 获取全部匹配沙箱，每页 20 条
```
$ qshell sandbox list --all --limit 20 -s running,paused -f json
```

6. 使用上一页游标继续
```
$ qshell sandbox list --next-token CURSOR -f json
```
