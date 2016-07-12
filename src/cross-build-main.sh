DIR=$(cd ../; pwd)
export GOPATH=$GOPATH:$DIR
GOOS=linux   GOARCH=386    go build -o ../bin/qshell_linux_386     main.go
GOOS=linux   GOARCH=amd64  go build -o ../bin/qshell_linux_amd64   main.go
GOOS=linux   GOARCH=arm    go build -o ../bin/qshell_linux_arm     main.go
GOOS=windows GOARCH=386    go build -o ../bin/qshell_windows_386.exe   main.go
GOOS=windows GOARCH=amd64  go build -o ../bin/qshell_windows_amd64.exe main.go
GOOS=darwin  GOARCH=386    go build -o ../bin/qshell_darwin_386    main.go
GOOS=darwin  GOARCH=amd64  go build -o ../bin/qshell_darwin_amd64  main.go
