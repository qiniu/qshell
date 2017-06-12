DIR=$(cd ../; pwd)
export GOPATH=$DIR:$GOPATH
go build main.go
