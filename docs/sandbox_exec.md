# 简介
`sandbox exec`（别名 `ex`）在沙箱中执行命令。默认前台模式会实时流式输出 stdout/stderr 并传递退出码；后台模式会打印进程 PID 后立即返回。

# 格式
```
qshell sandbox exec <sandboxID> -- <command...>
qshell sbx ex <sandboxID> -- <command...>
```

# 帮助文档
```
$ qshell sandbox exec -h
$ qshell sandbox exec --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量，或在当前目录 `.env` 文件中设置。

# 参数
- `sandboxID`：目标沙箱 ID
- `command`：要执行的命令（`--` 之后的所有参数）
- `-b, --background`：在后台运行命令，打印 PID 后返回
- `-c, --cwd`：设置命令的工作目录
- `-u, --user`：设置执行命令的用户
- `-e, --env`：设置环境变量（KEY=VALUE 格式，可多次指定）

# 示例
1. 在沙箱中执行命令
```
$ qshell sandbox exec sb-xxxxxxxxxxxx -- ls -la
$ qshell sbx ex sb-xxxxxxxxxxxx -- ls -la
```

2. 执行带管道的复杂命令
```
$ qshell sandbox exec sb-xxxxxxxxxxxx -- sh -lc 'cat /etc/os-release | head -5'
```

3. 将本地标准输入传给沙箱命令
```
$ echo "hello world" | qshell sbx ex sb-xxxxxxxxxxxx -- cat
$ cat file.txt | qshell sbx ex sb-xxxxxxxxxxxx -- wc -l
```

4. 后台运行命令
```
$ qshell sandbox exec sb-xxxxxxxxxxxx -b -- python server.py
$ qshell sbx ex sb-xxxxxxxxxxxx -b -- python server.py
```

5. 指定工作目录和用户
```
$ qshell sandbox exec sb-xxxxxxxxxxxx -c /app -u root -- npm install
```

6. 设置环境变量
```
$ qshell sandbox exec sb-xxxxxxxxxxxx -e PORT=3000 -e NODE_ENV=production -- node app.js
```
