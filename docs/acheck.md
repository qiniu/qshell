# 简介
`acheck` 查询异步抓取任务状态。

参考文档：[查询状态 (async fetch)](https://developer.qiniu.com/kodo/api/4097/asynch-fetch)

# 格式
```
 qshell acheck <Bucket> <ID> [flags]
```

# 参数
- Bucket：空间名，可以为公开空间或私有空间。 【必选】
- ID：异步 fetch 返回的任务 ID。 【必选】

详细的选项介绍，请参考：[查询状态 (async fetch)](https://developer.qiniu.com/kodo/api/4097/asynch-fetch)

# 示例
异步 fetch 返回 ID `eyJ6b25lIjoibmEwIiwicXVldWUiOiJTSVNZUEhVUy1KT0JTLVYzIiwicGFydF9pZCI6OSwib2Zmc2V0Ijo1NTEzMTU3fQ==`, bucket 名为 test：
```
qshell acheck test eyJ6b25lIjoibmEwIiwicXVldWUiOiJTSVNZUEhVUy1KT0JTLVYzIiwicGFydF9pZCI6OSwib2Zmc2V0Ijo1NTEzMTU3fQ==
```

