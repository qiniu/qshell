# 简介
`share-ls` 命令用来列举分享的目录和文件，对于每个文件，将会输出它的大小，文件类型和修改时间。

# 格式
```
qshell share-ls <Link>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell share-ls -h

// 详细文档（此文档）
$ qshell share-ls --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Link: 分享目录链接，或是通过 `create-share` 命令的 `--output` 选项输出的文件路径【必选】

# 选项
- --extract-code: 提取码，只能包含六位大小写字母或者数字，如果不填写，且 `Link` 不是通过 `create-share` 命令的 `--output` 选项输出的文件路径，将会用交互模式提示输入。【可选】
- --prefix: 列举前缀，如果不填写，将列举被分享目录下的所有文件。【可选】
- --limit: 限制列举条目的数量，当设置此选项，会在结束时输出 `Marker` 以供下次执行该命令时作为 `--marker` 选项的值，默认输出所有条目。【可选】
- --marker: 标记列举过程中的位置。【可选】

# 示例
```
$ qshell share-ls 'http://portal.qiniu.com/kodo-shares/verify?id=AGQEKDRxBBjbGmsKduQS9oFx59rz&token=qhtbC5YmDCO-WiPriuoCG_t4hZ1LboSOtRYSJXo_%3A9uJY8FiNrKjNrt4MpBx547jlgwr8aes15z5i8VY6l5SU6ga2IKWDBSGTv1jo-rOocklE7QqApzG6okJktZ36umLoqv9x1kuo5fNmgasLXowyTuHIM3kXsaV_DoXmvQsGr5ol6j4RtrmLcKdtXhpkGH8MfSjEgRV91Bx_Q_mSwpJ1028p8yZCSad_QOu_kSPxzeLZmWlUpAtO2oEXdbMTBxhTCH_3awCgqkgoogi0FQGP4zHxeFr0n3vj69DpmWqe6DiYbYLivCuU0kOF5Khv4I6-w6vjjdY'
Input Extract Code:
wp7gqc
go.mod   912     STANDARD        2024-09-04 02:21:13 +0000 UTC
go.sum   9092    STANDARD        2024-09-04 02:21:13 +0000 UTC
main     9093344 STANDARD        2024-09-04 02:21:14 +0000 UTC
pfop-example/
pfop-example/main.go    1094    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-download/
upload-download/main.go 4110    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-example/
upload-example/main.go  1486    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-lots-of-files/
upload-lots-of-files/main.go    2843    STANDARD        2024-09-04 02:21:13 +0000 UTC
Total size: 8.69MB
Folder number: 7
File number: 4

$ qshell share-ls share.json
go.mod   912     STANDARD        2024-09-04 02:21:13 +0000 UTC
go.sum   9092    STANDARD        2024-09-04 02:21:13 +0000 UTC
main     9093344 STANDARD        2024-09-04 02:21:14 +0000 UTC
pfop-example/
pfop-example/main.go    1094    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-download/
upload-download/main.go 4110    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-example/
upload-example/main.go  1486    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-lots-of-files/
upload-lots-of-files/main.go    2843    STANDARD        2024-09-04 02:21:13 +0000 UTC
Total size: 8.69MB
Folder number: 7
File number: 4

$ qshell share-ls share.json --limit=7
go.mod   912     STANDARD        2024-09-04 02:21:13 +0000 UTC
go.sum   9092    STANDARD        2024-09-04 02:21:13 +0000 UTC
main     9093344 STANDARD        2024-09-04 02:21:14 +0000 UTC
pfop-example/
pfop-example/main.go    1094    STANDARD        2024-09-04 02:21:13 +0000 UTC
upload-download/
upload-download/main.go 4110    STANDARD        2024-09-04 02:21:13 +0000 UTC
Marker: eyJjIjowLCJrIjoicHJlZml4L3VwbG9hZC1leGFtcGxlLyJ9
Total size: 8.68MB
Folder number: 5
File number: 2
```
