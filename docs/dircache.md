# 简介

`dircache`用来为指定的本地目录生成一个该目录下面所有的文件的列表，文件列表的每行分别为每个文件的相对路径，文件大小和最后修改时间，该文件列表可以作为`qupload`命令的`file_list`的参数。你可以修改该命令生成的文件，删除一些行，然后给`qupload`的参数`file_list`，这样就可以只上传指定的文件列表。

# 格式

```
qshell dircache <DirCacheRootPath> <DirCacheResultFile>
```

# 参数

|参数名|描述|
|------|------|
|DirCacheRootPath|本地需要对其做快照的路径，最好是绝对路径，比如`/Users/jemy/Demo1`这样的路径|
|DirCacheResultFile|快照结果的保存文件，可以是绝对路径或者相对路径|

# 示例

比如，要获取`/Users/jemy/Temp4`目录下面的文件列表，则使用

```
qshell dircache /Users/jemy/Temp4 temp4.list.txt
```

其中`temp4.list.txt`是你保存列表结果的文件。列举的结果以如下格式组织：

```
文件相对于<DirCacheRootPath>的相对路径\t文件大小(单位字节)\t文件上次修改时间(单位100纳秒)
```

比如这样的：

```
rk_video_not_play.mp4	3985210	14206026340000000
rtl1.flv	10342916	14205959890000000
sync_demo/array_enumeration.png	5262899	13953255140000000
sync_demo/demo2.gif	2685960	13966636230000000
sync_demo/golang.png	149366	14010291080000000
```