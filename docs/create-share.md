# 简介
`create-share` 命令为需要分享的目录创建授权链接。

# 格式
```
qshell create-share [kodo://]<Bucket>/<Prefix>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell create-share -h

// 详细文档（此文档）
$ qshell create-share --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket: 空间名称【必选】
- Prefix: 前缀【必选】

# 选项
- --extract-code: 提取码，只能包含六位大小写字母或者数字，如果不填写，将会自动生成。【可选】
- --validity-period: 有效时间，如果不填写，默认为 15 分钟。【可选】
- --output: 保存路径，以 JSON 格式保存输出内容，如果不填写，则直接以文本形式输出。【可选】

# 示例
```
$ qshell create-share kodo://bucketname/prefix
Link:
http://portal.qiniu.com/kodo-shares/verify?id=AGQEKDRxBBjbGmsKduQS9oFx59rz&token=qhtbC5YmDCO-WiPriuoCG_t4hZ1LboSOtRYSJXo_%3A9uJY8FiNrKjNrt4MpBx547jlgwr8aes15z5i8VY6l5SU6ga2IKWDBSGTv1jo-rOocklE7QqApzG6okJktZ36umLoqv9x1kuo5fNmgasLXowyTuHIM3kXsaV_DoXmvQsGr5ol6j4RtrmLcKdtXhpkGH8MfSjEgRV91Bx_Q_mSwpJ1028p8yZCSad_QOu_kSPxzeLZmWlUpAtO2oEXdbMTBxhTCH_3awCgqkgoogi0FQGP4zHxeFr0n3vj69DpmWqe6DiYbYLivCuU0kOF5Khv4I6-w6vjjdY
Extract Code:
wp7gqc
Expire:
2024-10-09 10:44:41 +0000

$ qshell create-share --output=share.json kodo://bucketname/prefix
```
