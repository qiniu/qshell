WIN86=qshell_windows_x86.exe
WIN64=qshell_windows_x64.exe
DARWIN=qshell_darwin_x64
LINUX86=qshell_linux_x86
LINUX64=qshell_linux_x64
LINUXARM=qshell_linux_arm

install:
	GOOS=darwin GOARCH=amd64 go build -o $(DARWIN) .
	cp ./$(DARWIN) /usr/local/bin/qshell && rm ./$(DARWIN)

all: linux windows arm

linux:
	GOOS=linux GOARCH=386 go build -o $(LINUX86) .
	GOOS=linux GOARCH=amd64 go build -o $(LINUX64) .

windows:
	GOOS=windows GOARCH=386 go build -o $(WIN86) .
	GOOS=windows GOARCH=amd64 go build -o $(WIN64) .

arm:
	GOOS=linux GOARCH=arm go build -o $(LINUXARM)
