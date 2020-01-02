VERSION=v2.4.1
WIN86=qshell-windows-x86-$(VERSION).exe
WIN64=qshell-windows-x64-$(VERSION).exe
DARWIN=qshell-darwin-x64-$(VERSION)
LINUX86=qshell-linux-x86-$(VERSION)
LINUX64=qshell-linux-x64-$(VERSION)
LINUXARM=qshell-linux-arm-$(VERSION)
LDFLAGS='-extldflags "-static"'
GO=GO111MODULE=on go

all: linux windows arm darwin

darwin:
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(DARWIN)
	@zip $(DARWIN).zip $(DARWIN)

linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GO) build -ldflags  $(LDFLAGS) -o $(LINUX86) .
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUX64) .
	@zip $(LINUX86).zip $(LINUX86)
	@zip $(LINUX64).zip $(LINUX64)


windows:
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GO) build -ldflags $(LDFLAGS) -o $(WIN86) .
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(WIN64) .
	@zip $(WIN86).zip $(WIN86)
	@zip $(WIN64).zip $(WIN64)

arm:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GO) build -ldflags $(LDFLAGS) -o $(LINUXARM)
	@zip $(LINUXARM).zip $(LINUXARM)

cleanzip:
	@rm $(LINUX86).zip $(LINUX64).zip $(LINUXARM).zip 2> /dev/null || true
	@rm $(WIN86).zip $(WIN64).zip 2> /dev/null || true
	@rm $(DARWIN).zip 2> /dev/null || true

cleanbin:
	@rm $(LINUX86) $(LINUX64) $(LINUXARM) $(DARWIN) $(WIN86) $(WIN64) 2> /dev/null || true

clean: cleanzip	cleanbin

upload:
	qshell rput devtools $(LINUX86).zip $(LINUX86).zip
	qshell rput devtools $(LINUX64).zip $(LINUX64).zip
	qshell rput devtools $(LINUXARM).zip $(LINUXARM).zip
	qshell rput devtools $(WIN86).zip $(WIN86).zip
	qshell rput devtools $(WIN64).zip $(WIN64).zip
	qshell rput devtools $(DARWIN).zip $(DARWIN).zip
