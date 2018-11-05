install:
	GOOS=darwin GOARCH=amd64 go build -o qshell_darwin_amd64 .
	cp ./qshell_darwin_amd64 /usr/local/bin/qshell

all:
	GOOS=windows GOARCH=386 go build -o qshell_windows_i386 .
	GOOS=windows GOARCH=amd64 go build -o qshell_windows_amd64 .
	GOOS=darwin GOARCH=amd64 go build -o qshell_darwin_amd64 .
	GOOS=linux GOARCH=386 go build -o qshell_linux_i386 .
	GOOS=linux GOARCH=amd64 go build -o qshell_linux_amd64 .

linux:
	GOOS=linux GOARCH=386 go build -o qshell_linux_i386 .
	GOOS=linux GOARCH=amd64 go build -o qshell_linux_amd64 .

windows:
	GOOS=windows GOARCH=386 go build -o qshell_windows_i386 .
	GOOS=windows GOARCH=amd64 go build -o qshell_windows_amd64 .
