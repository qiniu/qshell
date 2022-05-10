# 简介
`m3u8replace` 命令用来修改或删除七牛空间中 m3u8 播放列表中引用的切片路径中的域名。

# 格式
```
qshell m3u8replace <Bucket> <M3u8Key> [<NewDomain>]
``` 

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数
- Bucket：m3u8 文件所在空间，可以为公开空间或私有空间 【必选】
- M3u8Key：m3u8 文件的名字 【必选】
- NewDomain：引用切片的域名，如果不指定的话，则 m3u8 文件中引用切片使用相对路径，效果等同于转码时指定 `noDomain/1` 【可选】

# 示例
1 清除 m3u8 播放列表中切片引用路径中的域名，等同于转码时指定 `noDomain/1`
```
qshell m3u8replace if-pbl qiniu.m3u8
```

2 替换 m3u8 播放列表中切片引用路径中的域名，把旧的换成新的。
```
qshell m3u8replace if-pbl qiniu.m3u8 http://hls.example.com
```
 