DIR=$(cd ../; pwd)
export GOPATH=$GOPATH:$DIR
GOOS=linux   GOARCH=386    go build -o ../bin/qiniu_upload_linux_386         qiniu-upload.go
GOOS=linux   GOARCH=amd64  go build -o ../bin/qiniu_upload_linux_amd64       qiniu-upload.go
GOOS=windows GOARCH=386    go build -o ../bin/qiniu_upload_windows_386.exe   qiniu-upload.go
GOOS=windows GOARCH=amd64  go build -o ../bin/qiniu_upload_windows_amd64.exe qiniu-upload.go
GOOS=darwin  GOARCH=386    go build -o ../bin/qiniu_upload_darwin_386        qiniu-upload.go
GOOS=darwin  GOARCH=amd64  go build -o ../bin/qiniu_upload_darwin_amd64      qiniu-upload.go
