# 简介
`sandbox snapshot list` 查询一页持久快照。

# 使用
```shell
qshell sandbox snapshot list --sandbox-id <sandboxID> --name demo --limit 20 --format json
qshell sandbox snapshot list --next-token <cursor> --format json
```

名称和来源沙箱过滤参数均可省略。JSON 输出包含 Items 与 NextToken，游标为空表示没有下一页。翻页时保持相同过滤条件。
