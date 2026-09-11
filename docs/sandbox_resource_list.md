# 简介
`sandbox resource list`（别名 `ls`）查询指定沙箱已挂载的资源。

# 格式
```
qshell sandbox resource list <sandboxID> [--format <pretty|json>]
qshell sbx resource ls <sandboxID> [--format <pretty|json>]
```

# 参数
- `sandboxID`：沙箱 ID。
- `--format`：输出格式，支持 `pretty`（默认）或 `json`。

# 返回内容
Git 仓库资源包含资源 ID、仓库类型、仓库 URL 和挂载路径；Kodo 资源包含资源 ID、bucket、前缀、挂载路径和只读状态。服务端会对访问密钥和授权令牌做脱敏处理，CLI 不会输出这些敏感字段。

# 示例
```
$ qshell sandbox resource list sb-xxxxxxxxxxxx
$ qshell sbx resource ls sb-xxxxxxxxxxxx --format json
```
