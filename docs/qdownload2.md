# 简介
`qdownload2` 功能和 `qdownload` 一致，不过 `qdownload2` 通过命令行的方式来指定各个需要的参数，例如：
```
qshell qdownload2 --dest-dir=/home/jemy/temp --bucket=test
```

其所支持的命令参数列表，可以通过 `-h` 选项获得，参数含义参考：[qdownload](qdownload.md)
```
qshell qdownload2 -h                                         
```

```
Usage:
  qshell qdownload2 [-c <ThreadCount>]  [flags]

Flags:
      --bucket string                              storage bucket
      --cdn-domain string                          set the CDN domain name for downloading, the default is empty, which means downloading from the storage source site
      --check-hash                                 whether to verify the hash, if it is enabled, it may take a long time
      --dest-dir string                            local storage path, full path. default current dir
      --enable-slice --slice-file-size-threshold   whether to enable slice download, you need to pay attention to the configuration of --slice-file-size-threshold slice threshold option. Only when slice download is enabled and the size of the downloaded file is greater than the slice threshold will the slice download be started
      --get-file-api                               public storage cloud not support, private storage cloud support when has getfile api.
  -h, --help                                       help for qdownload2
      --io-host string                             io host of request
      --key-file string                            configure a file and specify the keys to be downloaded; if not configured, download all the files in the bucket
      --log-file string                            the output file of the download log is output to the file specified by record_root by default, and the specific file path can be seen in the terminal output (default "debug")
      --log-level string                           download log output level, optional values are debug,info,warn and error (default "debug")
      --log-rotate int                             the switching period of the download log file, the unit is day, (default 7)
      --prefix string                              only download files with the specified prefix
      --public                                     whether the space is a public space
      --record-root qshell                         path to save download record information, including log files and download progress files; the default is qshell download directory
      --referer string                             if the CDN domain name is configured with domain name whitelist anti-leech, you need to specify a referer address that allows access
      --remove-temp-while-error                    when the download encounters an error, delete the previously downloaded part of the file cache
      --save-path-handler string                   specify a callback function; when constructing the save path of the file, this option is preferred for construction. If not configured, $dest_dir + $ file separator + $Key will be used for construction. This function is implemented through the template of the Go language. The func command is used for function verification. For the specific syntax, please refer to the description of the func command.
      --slice-concurrent-count int                 concurrency of slice downloads (default 10)
      --slice-file-size-threshold int              file threshold for downloading slices. When slice downloading is enabled and the file size is greater than this threshold, slice downloading will be enabled; unit:B (default 41943040)
      --slice-size int                             slice size; when using slice download, the size of each slice; unit:B (default 4194304)
      --suffixes string                            only download files with the specified suffixes
  -c, --thread int                                 num of threads to download files (default 5)
```