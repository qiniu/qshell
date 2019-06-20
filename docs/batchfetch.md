# 简介

`batchfetch`命令用来批量抓取远程地址到七牛存储空间。


# 格式

```
qshell batchfetch [-F <Delimiter>] [--success-list <SuccessFileName>] [--failure-list <failureFileName>] [-i <FetchUrlsFile>] <Bucket>
```

# 帮助
```
qshell batchcopy -h
```

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey`和`Name`的情况下使用。

# 选项和参数

| 参数名           | 描述                                                                                          |
|------------------|-----------------------------------------------------------------------------------------------|
| success-list选项, 接受一个文件名参数 | 导出抓取成功的url到该文件， 文件有两列, 第一列是抓取的URL, 第二列是保存在存储空间中的文件名字 |
| failure-list选项, 接受一个文件名参数 | 导出抓取失败的url到该文件， 文件有两列, 第一列是抓取的URL， 第二列是抓取报错信息| 
| i 选项， 接受一个文件名<FetchUrlsFile>参数 | 该文件是要抓取的文件URL地址和要保存在七牛的存储空间中的名称, 如果该选项没有指定，默认从标准输入读取内容 |

<FetchUrlsFile>文件内容的有如下几种格式：

**模式一:**

```
文件链接1<Delimiter>保存名称1
文件链接2<Delimiter>保存名称2
文件链接3<Delimiter>保存名称3
...
```
其中<Delimiter> 表示分隔符, 默认使用空白进行分隔（空格，\t, \n), 如果要抓取的地址或者保存的文件名中有空格， 可以使用-F选项指定分隔符

例如：


```
http://img.abc.com/0/000/484/0000484193.fid	2009-10-14/2922168_b.jpg
http://img.abc.com/0/000/553/0000553777.fid	2009-07-01/2270194_b.jpg
http://img.abc.com/0/000/563/0000563511.fid	2009-03-01/1650739_s.jpg
http://img.abc.com/0/000/563/0000563514.fid	2009-05-01/1953696_m.jpg
http://img.abc.com/0/000/563/0000563515.fid	2009-02-01/1516376_s.jpg
```
上面的方式最终抓取保存在空间中的文件名字是：


```
2009-10-14/2922168_b.jpg
2009-07-01/2270194_b.jpg
2009-03-01/1650739_s.jpg
2009-05-01/1953696_m.jpg
2009-02-01/1516376_s.jpg
```

**模式二:**

```
文件链接1
文件链接2
文件链接3
...
```

上面的方式也是支持的，这种方式的情况下，文件保存的名字将从指定的文件链接里面自动解析。

例如：

```
http://img.abc.com/0/000/484/0000484193.fid
http://img.abc.com/0/000/553/0000553777.fid
http://img.abc.com/0/000/563/0000563511.fid
http://img.abc.com/0/000/563/0000563514.fid
http://img.abc.com/0/000/563/0000563515.fid
```

其抓取后保存在空间中的文件名字是：


```
0/000/484/0000484193.fid
0/000/553/0000553777.fid
0/000/563/0000563511.fid
0/000/563/0000563514.fid
0/000/563/0000563515.fid
```

# 使用示例

假如我们的AccessKey="test-ak", SecretKey="test-sk", 我给自己账号起了个名字Name="myself"

第一步:
检查qshell本地数据库有没有该账号，如果有该账号，会打印出来该账号的信息

```
$ qshell user lookup myself
```

如果有该账号，可以使用

```
$ qshell user cu myself
```
切换到该账号, 如果您配置了自动补全（配置方法参考README.md)， 在命令行输入
```
$ qshell user cu <TAB>
```
会自动补全本地数据库的账户名字

如果没有该账号，需要使用qshell account 添加账号到qshell的本地数据库, 其中<Your AccountName>可以自定义, 改名字的作用只是用来在本地数据库中唯一表示账户名称

```
$ qshell account <Your AccessKey> <Your SecretKey> <Your AccountName>
```

第二步:

使用batchfetch命令操作, 假如我要操作的bucket="test-bucket", 要预取的文件地址列表保存在文件batchfetchurls.txt：

```
$ qshell batchfetch test-bucket -i batchfetchurls.txt
```

如果想导出fetch成功，失败的列表分别到文件fetch_success.txt, fetch_failure.txt，可以使用如下命令:

```
$ qshell batchfetch test-bucket -i batchfetchurls.txt --success-list fetch_success.txt --failure-list fetch_failure.txt
```
