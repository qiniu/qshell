DIR=$(cd ../; pwd)
export GOPATH=$DIR:$GOPATH
GOOS=windows GOARCH=386    go build -o ../bin/qshell_windows_x86.exe   main.go
GOOS=windows GOARCH=amd64  go build -o ../bin/qshell_windows_x64.exe main.go
GOOS=darwin  GOARCH=386    go build -o ../bin/qshell_darwin_x86    main.go
GOOS=darwin  GOARCH=amd64  go build -o ../bin/qshell_darwin_x64  main.go
GOOS=linux  GOARCH=386    go build -o ../bin/qshell_linux_x86    main.go
GOOS=linux  GOARCH=amd64  go build -o ../bin/qshell_linux_x64  main.go
