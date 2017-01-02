# 简介

`prefop` 命令用来根[pfop请求](http://developer.qiniu.com/code/v6/api/dora-api/pfop/pfop.html)得到的 `PersistentId` 来查询七牛数据处理的状态或结果。

参考文档：[持久化处理状态查询 (prefop)](http://developer.qiniu.com/code/v6/api/dora-api/pfop/prefop.html)

# 格式

```
qshell prefop <PersistentId>
```

# 鉴权
无

# 参数

|参数名|描述|
|------|---------|
|PersistentId|持久化处理的Id|

# 示例

查询 `z0.58632a1945a2650cfd5fc8b1` 对应的持久化处理结果：

```
qshell prefop z0.58632a1945a2650cfd5fc8b1
```

输出

```
Id: z0.58632a1945a2650cfd5fc8b1
Code: 0
Desc: The fop was completed successfully
InputBucket: video
InputKey: bjsp/c70f228f-5133-a2cc-a811-5dfa5433996f.mp4
Pipeline: 1380710990.jcsp_bjsp
Reqid: SX4AAFKzn5IZTJQU

	Cmd:	avthumb/mp4|saveas/dmlkZW86YmpzcC9jNzBmMjI4Zi01MTMzLWEyY2MtYTgxMS01ZGZhNTQzMzk5NmYubXA0
	Code:	0
	Desc:	The fop was completed successfully
	Hash:	lrttDVJrZJYVqV46PgfJojip6Zrn
	Key:	bjsp/c70f228f-5133-a2cc-a811-5dfa5433996f.mp4
```
