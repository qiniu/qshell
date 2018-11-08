install:
	GOOS=darwin GOARCH=amd64 go build -o qshell_darwin_amd64 .
	cp ./qshell_darwin_amd64 /usr/local/bin/qshell && rm ./qshell_darwin_amd64

all:
	GOOS=windows GOARCH=386 go build -o qshell_windows_x86.exe .
	GOOS=windows GOARCH=amd64 go build -o qshell_windows_x64.exe .
	GOOS=darwin GOARCH=amd64 go build -o qshell_darwin_x64 .
	GOOS=linux GOARCH=386 go build -o qshell_linux_x86 .
	GOOS=linux GOARCH=amd64 go build -o qshell_linux_x64 .

linux:
	GOOS=linux GOARCH=386 go build -o qshell_linux_x86 .
	GOOS=linux GOARCH=amd64 go build -o qshell_linux_x64 .

windows:
	GOOS=windows GOARCH=386 go build -o qshell_windows_x86 .
	GOOS=windows GOARCH=amd64 go build -o qshell_windows_x64 .
