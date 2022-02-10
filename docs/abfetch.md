# 简介
`abfetch` 使用异步抓取接口抓取网络资源到七牛存储空间。

参考文档：[异步抓取 (async fetch)](https://developer.qiniu.com/kodo/api/4097/asynch-fetch)

# 格式
```
qshell abfetch [-i <URLList>][-b <CallbackBody>][-T <CallbackHost>][-a <CallbackUrl>][-e <FailureList>][-t <DownloadHostHeader>][-g <StorageType>][-s <SuccessList>][-c <ThreadCount>] <Bucket>
```

# 选项
| 选项 | 说明                                                         |
| -i   | 要抓取的资源列表， 一行一个资源                              |
| -b   | 回调的http Body                                              |
| -T   | 回调时的HOST 头                                              |
| -a   | 回调的请求地址                                               |
| -e   | 抓取失败导出的文件列表                                       |
| -t   | 下载资源时使用的HOST头                                       |
| -g   | 抓取的资源存储在七牛存储空间的类型， 分为低频存储， 标准存储 |
| -s   | 抓取成功后导出到的文件                                       |
| -c   | 抓取指定的线程数目                                                    |

详细的选项介绍，请参考：[异步抓取 (async fetch)](https://developer.qiniu.com/kodo/api/4097/asynch-fetch)


# 参数
<Bucket> 为七牛存储空间名


# 例子
假如我有3个资源要抓取，地址分别为：
http://test.com/test1.txt
http://test.com/test2.txt
http://test.com/test3.txt

需要抓取这三个资源保存在七牛存储空间"test"中

#### 第一步：
在当前目录下创建名为"urls.txt"的文件， 文件内容为

```
http://test.com/test1.txt
http://test.com/test2.txt
http://test.com/test3.txt
```
每行一个地址 。

#### 第二步:
使用如下的命令就可以抓取资源到存储"test"中
```
$ qshell abfetch -i urls.txt test
```

但是这样我们不知道哪些成功了，哪些抓取失败了，可以使用选项-e 导出失败列表到文件"failure.txt"中:
```
$ qshell abfetch -i urls.txt -e failure.txt test
```

如果要提高请求的并发量， 可以使用选项-c 指定提交的线程数, 下面的命令指定线程数为100:
```
$ qshell abfetch -i urls.txt -e failure.txt -c 100 test
```
线程数只能决定请求接口后台服务器的快慢， 提交的请求会到服务器处理队列中， 如果队列中有很多要抓取的资源，抓取速度不一定会提高，所以适当设置线程数

# 文件大小
异步接口暂时没办法判断是否抓取成功， 当异步接口返回的数据wait是-1时，表示抓取过这个文件，这时程序会用stat接口去存储获取文件的信息，如果可以获取到，说明抓取成功了；如果wait为-1， 重试三次stat都失败，那么认为抓取失败。

获取异步接口的返回结果是通过轮训进行的， 每次轮训的时间间隔取决于文件的大小， 大的文件时间间隔长点，小的文件时间间隔短点。
因此可以在抓取的资源文件后面附上文件大小来帮助程序估算大概的时间间隔。

比如要抓取的资源地址为：
http://test.com/test1.txt
http://test.com/test2.txt
http://test.com/test3.txt

test1.txt 大小为400字节
test2.txt 大小为300340字节
test3.txt 大小为22039字节

假如把这些信息保存在"fetch_input.txt"文件中, 文件内容格式为：
http://test.com/test1.txt    400
http://test.com/test2.txt    300340
http://test.com/test3.txt    22039

程序会根据文件大小估算异步抓取的时间，时间到后，会访问异步结果获得处理结果。如果没有指定文件大小，默认会有轮训时间。
