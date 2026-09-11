# 简介
`sandbox resource update`（别名 `up`、`update-token`）按资源 ID 更新 Git 仓库资源的授权令牌。对正在运行的沙箱，新令牌会立即应用到对应资源。

# 格式
```
qshell sandbox resource update <sandboxID> <resourceID> --token <token>
qshell sbx resource up <sandboxID> <resourceID> -t <token>
```

# 参数
- `sandboxID`：沙箱 ID。
- `resourceID`：Git 仓库资源 ID，可通过 `sandbox resource list` 查询。
- `-t, --token`：新的 Git 仓库授权令牌，必填。
- `--authorization-token`：`--token` 的长参数别名。

# 示例
```
$ qshell sandbox resource update sb-xxxxxxxxxxxx res-xxxxxxxxxxxx --token ghp-xxx
$ qshell sbx resource up sb-xxxxxxxxxxxx res-xxxxxxxxxxxx -t ghp-xxx
```
