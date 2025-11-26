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
      --accelerate                       enable uploading acceleration
      --bucket string                    bucket
      --callback-body string             upload callback body
  -T, --callback-host string             upload callback host
  -l, --callback-urls string             upload callback urls, separated by comma
      --check-exists                     check file key whether in bucket before upload
      --check-hash                       check hash
      --check-size                       check file size
      --detect-mime int                  Turn on the MimeType detection function and perform detection according to the following rules; if the correct value cannot be detected, application/octet-stream will be used by default.
                                         If set to a value of 1, the file MimeType information passed by the uploader will be ignored, and the MimeType value will be detected in the following order:
                                         	1. Detection content;
                                         	2. Check the file extension;
                                         	3. Check the Key extension.
                                         The default value is set to 0. If the uploader specifies MimeType (except application/octet-stream), this value will be used directly. Otherwise, the MimeType value will be detected in the following order:
                                         	1. Check the file extension;
                                         	2. Check the Key extension;
                                         	3. Detect content.
                                         Set to a value of -1 and use this value regardless of what value is specified on the uploader.
      --end-user string                  Owner identification
  -e, --failure-list string              upload failure file list
      --file-list string                 file list to upload
      --file-type int                    set storage type of file, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage, 4:ARCHIVE_IR storage, 5:INTELLIGENT_TIERING
  -h, --help                             help for qupload2
      --ignore-dir                       ignore the dir in the dest file key
      --key-prefix string                key prefix prepended to dest file key
      --log-file string                  log file
      --log-level string                 log level (default "debug")
      --log-rotate int                   log rotate days (default 7)
      --overwrite                        overwrite the file of same key in bucket
  -w, --overwrite-list string            upload success (overwrite) file list
      --persistent-notify-url string     URL to receive notification of persistence processing results. It must be a valid URL that can make POST requests normally on the public Internet and respond successfully. The content obtained by this URL is consistent with the processing result of the persistence processing status query. To send a POST request whose body format is application/json, you need to read the body of the request in the form of a read stream to obtain it.
      --persistent-ops string            List of pre-transfer persistence processing instructions that are triggered after successful resource upload. This parameter is not supported when fileType=2 or 3 (upload archive storage or deep archive storage files). Supports magic variables and custom variables. Each directive is an API specification string, and multiple directives are separated by ;.
      --persistent-pipeline string       Transcoding queue name. After the resource is successfully uploaded, an independent queue is designated for transcoding when transcoding is triggered. If it is empty, it means that the public queue is used, and the processing speed is slower. It is recommended to use a dedicated queue.
      --put-threshold int                chunk upload threshold, unit: B (default 8388608)
      --record-root string               record root dir, and will save record info to the dir(db and log), default <UserRoot>/.qshell
      --rescan-local                     rescan local dir to upload newly add files
      --resumable-api-v2                 use resumable upload v2 APIs to upload
      --resumable-api-v2-part-size int   the part size when use resumable upload v2 APIs to upload (default 4194304)
      --sequential-read-file             File reading is sequential and does not involve skipping; when enabled, the uploading fragment data will be loaded into the memory. This option may increase file upload speed for mounted network filesystems.
      --skip-file-prefixes string        skip files with these file prefixes
      --skip-fixed-strings string        skip files with the fixed string in the name
      --skip-path-prefixes string        skip files with these relative path prefixes
      --skip-suffixes string             skip files with these suffixes
      --src-dir string                   src dir to upload
  -s, --success-list string              upload success file list
      --thread-count int                 multiple thread count (default 1)
      --traffic-limit uint               Upload request single link speed limit to control client bandwidth usage. The speed limit value range is 819200 ~ 838860800, and the unit is bit/s.
      --up-host string                   upload host
      --worker-count int                 the number of concurrently uploaded parts of a single file in resumable upload (default 3)

Global Flags:
      --colorful        console colorful mode
  -C, --config string   set config file (default is $HOME/.qshell.json)
  -D, --ddebug          deep debug mode
  -d, --debug           debug mode
      --doc             document of command
  -L, --local           use current directory qshell workspace (default is $HOME/.qshell)
      --silence         silence mode, The console only outputs warnings、errors and some important information
```
