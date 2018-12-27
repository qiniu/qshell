WIN86=qshell_windows_x86.exe
WIN64=qshell_windows_x64.exe
DARWIN=qshell_darwin_x64
LINUX86=qshell_linux_x86
LINUX64=qshell_linux_x64
LINUXARM=qshell_linux_arm
LDFLAGS='-extldflags "-static"'
GO=GO111MODULE=on go

all: linux windows arm darwin

darwin:
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(DARWIN)

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GO) build -ldflags  $(LDFLAGS) -o $(LINUX86) .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUX64) .

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GO) build -ldflags $(LDFLAGS) -o $(WIN86) .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(WIN64) .

arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GO) build -ldflags $(LDFLAGS) -o $(LINUXARM)
