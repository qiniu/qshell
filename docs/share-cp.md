# 简介
`share-cp` 命令用来从目录分享链接内下载单个文件或按指定前缀批量下载文件。

# 格式
```
qshell share-cp <Link> --to=<ToLocalPath>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell share-cp -h

// 详细文档（此文档）
$ qshell share-cp --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Link: 分享目录链接，或是通过 `create-share` 命令的 `--output` 选项输出的文件路径【必选】

# 选项
- --extract-code: 提取码，只能包含六位大小写字母或者数字，如果不填写，且 `Link` 不是通过 `create-share` 命令的 `--output` 选项输出的文件路径，将会用交互模式提示输入。【可选】
- --from: 下载共享目录内的子目录或是子文件，如果不填写，将列举被分享的目录。【可选】
- --to: 下载目标路径。【必填】

# 示例
```
$ qshell share-cp 'http://portal.qiniu.com/kodo-shares/verify?id=AGQEKDRxBBjbGmsKduQS9oFx59rz&token=qhtbC5YmDCO-WiPriuoCG_t4hZ1LboSOtRYSJXo_%3A9uJY8FiNrKjNrt4MpBx547jlgwr8aes15z5i8VY6l5SU6ga2IKWDBSGTv1jo-rOocklE7QqApzG6okJktZ36umLoqv9x1kuo5fNmgasLXowyTuHIM3kXsaV_DoXmvQsGr5ol6j4RtrmLcKdtXhpkGH8MfSjEgRV91Bx_Q_mSwpJ1028p8yZCSad_QOu_kSPxzeLZmWlUpAtO2oEXdbMTBxhTCH_3awCgqkgoogi0FQGP4zHxeFr0n3vj69DpmWqe6DiYbYLivCuU0kOF5Khv4I6-w6vjjdY' --to=shared/
Input Extract Code:
wp7gqc

$ qshell share-cp share.json --from=main.go --to=shared/
```
