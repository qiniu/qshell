# 简介
`qupload2` 功能和 `qupload` 一致，不过 `qupload2` 通过命令行的方式来指定各个需要的参数，例如：
```
qshell qupload2 --src-dir=/home/jemy/temp --bucket=test
```

其所支持的命令参数列表，可以通过 `-h` 选项获得，参数含义参考 `qupload` 命令类似的选项：[qupload](qupload.md)
例子：
`qupload2` 的 `--bucket` 选项含义可参考 `qupload` 的 `bucket` 配置；
`qupload2` 的 `--check-hash` 选项含义可参考 `qupload` 的 `check_hash` 配置；

```
jemy•~» qshell qupload2 -h                                         
```

```
Batch upload files to the qiniu bucket

Usage:
  qshell qupload2 [flags]

Flags:
      --bucket string                    bucket
  -T, --callback-host string             upload callback host
  -l, --callback-urls string             upload callback urls, separated by comma
      --check-exists                     check file key whether in bucket before upload
      --check-hash                       check hash
      --check-size                       check file size
  -e, --failure-list string              upload failure file list
      --file-list string                 file list to upload
      --file-type int                    set storage type of file, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage
  -h, --help                             help for qupload2
      --ignore-dir                       ignore the dir in the dest file key
      --key-prefix string                key prefix prepended to dest file key
      --log-file string                  log file
      --log-level string                 log level (default "debug")
      --log-rotate int                   log rotate days (default 7)
      --overwrite                        overwrite the file of same key in bucket
  -w, --overwrite-list string            upload success (overwrite) file list
      --put-threshold int                chunk upload threshold, unit: B (default 8388608)
      --record-root string               record root dir, and will save record info to the dir(db and log), default <UserRoot>/.qshell
      --rescan-local                     rescan local dir to upload newly add files
      --resumable-api-v2                 use resumable upload v2 APIs to upload
      --resumable-api-v2-part-size int   the part size when use resumable upload v2 APIs to upload (default 4194304)
      --skip-file-prefixes string        skip files with these file prefixes
      --skip-fixed-strings string        skip files with the fixed string in the name
      --skip-path-prefixes string        skip files with these relative path prefixes
      --skip-suffixes string             skip files with these suffixes
      --src-dir string                   src dir to upload
  -s, --success-list string              upload success file list
      --thread-count int                 multiple thread count (default 1)
      --up-host string                   upload host
      --worker-count int                 the number of concurrently uploaded parts of a single file in resumable upload (default 3)
```
