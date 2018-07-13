DIR=$(cd ../; pwd)
export GOPATH=$DIR:$GOPATH
GOOS=linux   GOARCH=amd64  go build -o qshell_linux_x64   main.go
GOOS=linux   GOARCH=386    go build -o qshell_linux_x86   main.go
GOOS=linux   GOARCH=arm    go build -o qshell_linux_arm   main.go