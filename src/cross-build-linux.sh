DIR=$(cd ../; pwd)
export GOPATH=$GOPATH:$DIR
GOOS=linux   GOARCH=amd64  go build -o qshell_linux_amd64   main.go
