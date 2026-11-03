# 简介
`sandbox snapshot create` 创建持久快照。

# 使用
```shell
qshell sandbox snapshot create <sandboxID> --idempotency-key <key> [--name demo:v1] [--memory=false] [--format json]
```

幂等键必填，调用前生成并保存（例如使用 uuidgen）。请求超时、断网或取消后，用原键、原沙箱 ID 和完全相同的请求体重试，不能直接换新键；取消等待不会取消已接受的捕获。

`--memory` 默认 true；false 仅保存文件系统，从快照创建沙箱时冷启动，进程和内存状态不会保留。`--name` 不能包含 `/`，命名空间由账号分配。
