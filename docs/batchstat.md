# 简介

`batchstat`命令用来批量查询七牛空间中文件的基本信息。

# 格式

```
qshell batchstat <Bucket> <KeyListFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|
|---------|-----------|
|Bucket|空间名，可以为公开空间或私有空间|
|KeyListFile|待查询文件名列表，按行分隔|


# 示例

- 我们将查询空间`7qiniu`中的一些文件的基本信息，待查询文件列表`listFile` 的内容为：

```
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000000.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000001.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000002.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000003.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000004.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000005.ts

```

- 使用如下命令进行批量查询

```
$ qshell batchstat 7qiniu listFile
```

- 输出 Key、Fsize、Hash、MimeType、PutTime 以`\t`分隔：

```
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000000.ts 92308   Fk8Uf2SHbQ4S2-cXHINuRc_rooNA    video/mp2t  15003760414606314
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000001.ts 91556   FpJP2nfipuLVc6QGvvcb868Rd0pO    video/mp2t  15003760414789673
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000002.ts 92496   FvBjZPch6cf52t2x0ZQBngqS1KTp    video/mp2t  15003760417159000
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000003.ts 92308   FoEgsbzdrcLuj_Fo5FeTI3w1jFHJ    video/mp2t  15003760419154144
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000004.ts 92308   FkYNctlf1JOGcJa-WzWgxsqcBjX6    video/mp2t  15003760422258065
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000005.ts 92120   Fh4Fwhu3dMUGbd3jE5OmRtfVZLv4    video/mp2t  15003760423842522
```

