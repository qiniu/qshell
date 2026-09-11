# 简介
`sandbox resource`（别名 `resources`）管理沙箱启动时挂载的资源。

# 格式
```
qshell sandbox resource
qshell sbx resource
```

# 子命令
- `list`：查询指定沙箱已挂载的资源，返回内容不包含访问密钥和授权令牌。
- `update`：按资源 ID 更新 Git 仓库资源的授权令牌。

# 示例
```
$ qshell sandbox resource list sb-xxxxxxxxxxxx
$ qshell sandbox resource list sb-xxxxxxxxxxxx --format json
$ qshell sandbox resource update sb-xxxxxxxxxxxx res-xxxxxxxxxxxx --token ghp-xxx
```
