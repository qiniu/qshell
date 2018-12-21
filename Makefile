WIN86=qshell_windows_x86.exe
WIN64=qshell_windows_x64.exe
DARWIN=qshell_darwin_x64
LINUX86=qshell_linux_x86
LINUX64=qshell_linux_x64
LINUXARM=qshell_linux_arm

all: linux windows arm darwin

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(DARWIN)

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags '-extldflags "-static"' -o $(LINUX86) .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static"' -o $(LINUX64) .

windows:
	GOOS=windows GOARCH=386 go build -o $(WIN86) .
	GOOS=windows GOARCH=amd64 go build -o $(WIN64) .

arm:
	GOOS=linux GOARCH=arm go build -o $(LINUXARM)
