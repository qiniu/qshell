# 简介
`saveas` 命令生成一个资源链接，该链接可以实现实时处理并且将处理结果存储到空间中。一般用在图片实时处理并同时持久化的过程中。

# 格式
```
qshell saveas <PublicUrlWithFop> <SaveBucket> <SaveKey>
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数
- PublicUrlWithFop：带实时处理指令的资源公开外链。 【必须】
- SaveBucket：处理结果保存的空间。 【必须】
- SaveKey：处理结果保存的文件名字。 【必须】

# 示例
1 我们需要对空间 `if-pbl` 里面的文件 `qiniu.png` 进行实时处理并且把结果保存在空间 `if-pbl `中，保存的文件名字为 `qiniu_1.jpg`。
我们可以用如下指令：
```
$ qshell saveas 'http://if-pbl.qiniudn.com/qiniu.png?imageView2/0/format/jpg' 'if-pbl' 'qiniu_1.jpg'
```

生成的结果外链：
```
http://if-pbl.qiniudn.com/qiniu.png?imageView2/0/format/jpg|saveas/aWYtcGJsOnFpbml1XzEuanBn/sign/TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:Rits4ikIlxTig5h0N3jAPbGdmmQ=
```

从上面的结果看，该命令自动为外链加上了saveas参数并且做了签名。

2 上面的例子是针对公开空间的，那么私有空间中的文件该如何处理呢？其实还是一样的。对于私有空间 `if-pri` 里面的文件 `qiniu.png`，我们一样按照上面的方法先生成公开的访问外链：
```
$ qshell saveas 'http://if-pri.qiniudn.com/qiniu.png?imageView2/0/format/jpg' 'if-pri' 'qiniu_1.jpg'
```

得到：
```
http://if-pri.qiniudn.com/qiniu.png?imageView2/0/format/jpg|saveas/aWYtcHJpOnFpbml1XzEuanBn/sign/TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:IM_TqyMu3rSRLuhgP3maTktRjPw=
```

但是上面的外链是无法直接访问的，我们还需要对这个外链进行私有空间访问的授权，使用`privateurl` 命令。
```
$ qshell privateurl 'http://if-pri.qiniudn.com/qiniu.png?imageView2/0/format/jpg|saveas/aWYtcHJpOnFpbml1XzEuanBn/sign/TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:IM_TqyMu3rSRLuhgP3maTktRjPw='
```

得到最终可以访问的链接：
```
http://if-pri.qiniudn.com/qiniu.png?imageView2/0/format/jpg|saveas/aWYtcHJpOnFpbml1XzEuanBn/sign/TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:IM_TqyMu3rSRLuhgP3maTktRjPw=&e=1430898125&token=TQt-iplt8zbK3LEHMjNYyhh6PzxkbelZFRMl10MM:nyLNxkJLSj2Z0-Ht-WIiISrMX1Y=
```