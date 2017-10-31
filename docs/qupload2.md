# 简介

`qupload2` 功能和 `qupload` 一致，不过 `qupload2` 通过命令行的方式来指定各个需要的参数，例如：

```
qshell qupload2 -src-dir=/home/jemy/temp -bucket=test
```

其所支持的命令参数列表，可以通过 `-h` 选项获得，参数含义参考：[qupload](qupload.md)

```
jemy•~» qshell qupload2 -h                                                                                                                                                                                                                                      
Usage of qupload2:
  -bind-nic-ip string
    	local network interface card to bind
  -bind-rs-ip string
    	rs host ip to bind
  -bind-up-ip string
    	upload host ip to bind
  -bucket string
    	bucket
  -check-exists
    	check file key whether in bucket before upload
  -check-hash
    	check hash
  -check-size
    	check file size
  -file-list string
    	file list to upload
  -ignore-dir
    	ignore the dir in the dest file key
  -key-prefix string
    	key prefix prepended to dest file key
  -log-file string
    	log file
  -log-level string
    	log level (default "info")
  -log-rotate int
    	log rotate days (default 1)
  -overwrite
    	overwrite the file of same key in bucket
  -put-threshold int
    	chunk upload threshold
  -rescan-local
    	rescan local dir to upload newly add files
  -skip-file-prefixes string
    	skip files with these file prefixes
  -skip-fixed-strings string
    	skip files with the fixed string in the name
  -skip-path-prefixes string
    	skip files with these relative path prefixes
  -skip-suffixes string
    	skip files with these suffixes
  -src-dir string
    	src dir to upload
  -thread-count int
    	multiple thread count
  -up-host string
    	upload host
```