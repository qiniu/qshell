# 简介

`ip`是一个查询ip归属地和运营商的工具，根据淘宝的API来获取结果。

API地址：[http://ip.taobao.com/service/getIpInfo.php](http://ip.taobao.com/service/getIpInfo.php)。该工具一次可以查询多个ip的信息。

# 格式

```
qshell ip <Ip1> [<Ip2> [<Ip3> ...]]]
```

# 参数

|参数名称|描述|可选|
|----|----------|---|
|Ip1|第一个ip地址|N|
|Ip2|第二个ip地址|Y|
|Ip3|第三个ip地址|Y|
|IpN|第N个ip地址|Y|

# 示例

查询`180.154.236.170`和`192.168.1.1`的归属地和运营商。

```
qshell ip 180.154.236.170  192.168.1.1
```

输出：

```
Ip: 180.154.236.170      => Code: 0	Country: 中国	Area: 华东	Region: 上海市	City: 上海市	County: 	Isp: 电信
Ip: 192.168.1.1          => Code: 0	Country: 未分配或者内网IP	Area: 	Region: 	City: 	County: 	Isp:
```