VERSION=$(if $(RELEASE_VERSION),$(RELEASE_VERSION),UNSTABLE)
$(info VERSION: $(VERSION))

DARWIN=qshell-$(VERSION)-darwin-x64
DARWINARM64=qshell-$(VERSION)-darwin-arm64
LINUX86=qshell-$(VERSION)-linux-x86
LINUX64=qshell-$(VERSION)-linux-x64
LINUXARM=qshell-$(VERSION)-linux-arm
LINUXARM64=qshell-$(VERSION)-linux-arm64
LINUXMIPS=qshell-$(VERSION)-linux-mips
LINUXMIPSLE=qshell-$(VERSION)-linux-mipsle
LINUXMIPS64=qshell-$(VERSION)-linux-mips64
LINUXMIPS64LE=qshell-$(VERSION)-linux-mips64le
LINUXLOONG64=qshell-$(VERSION)-linux-loong64
LINUXRISCV64=qshell-$(VERSION)-linux-riscv64
WIN86=qshell-$(VERSION)-windows-x86
WIN64=qshell-$(VERSION)-windows-x64
WINARM=qshell-$(VERSION)-windows-arm
WINARM64=qshell-$(VERSION)-windows-arm64

LDFLAGS='-X 'github.com/qiniu/qshell/v2/iqshell/common/version.version=$(VERSION)' -extldflags '-static''
GO=GO111MODULE=on go

all: linux windows darwin

.PHONY: linux windows darwin

darwin: $(DARWIN).zip $(DARWINARM64).zip

linux: $(LINUX86).zip $(LINUX64).zip $(LINUXARM).zip $(LINUXARM64).zip $(LINUXMIPS).zip $(LINUXMIPSLE).zip $(LINUXMIPS64).zip $(LINUXMIPS64LE).zip $(LINUXLOONG64).zip $(LINUXRISCV64).zip

windows: $(WIN86).zip $(WIN64).zip $(WINARM).zip $(WINARM64).zip

qshell-$(VERSION)-darwin-%.zip: qshell-$(VERSION)-darwin-%
	zip $@ $<
qshell-$(VERSION)-linux-%.zip: qshell-$(VERSION)-linux-%
	zip $@ $<
qshell-$(VERSION)-windows-%.zip: qshell-$(VERSION)-windows-%.exe
	zip $@ $<

$(DARWIN):
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(DARWIN) ./main/
$(DARWINARM64):
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags $(LDFLAGS) -o $(DARWINARM64) ./main/
$(LINUX86):
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GO) build -ldflags  $(LDFLAGS) -o $(LINUX86) ./main/
$(LINUX64):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUX64) ./main/
$(LINUXARM):
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GO) build -ldflags $(LDFLAGS) -o $(LINUXARM) ./main/
$(LINUXARM64):
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUXARM64) ./main/
$(LINUXMIPS):
	CGO_ENABLED=0 GOOS=linux GOARCH=mips $(GO) build -ldflags $(LDFLAGS) -o $(LINUXMIPS) ./main/
$(LINUXMIPSLE):
	CGO_ENABLED=0 GOOS=linux GOARCH=mipsle $(GO) build -ldflags $(LDFLAGS) -o $(LINUXMIPSLE) ./main/
$(LINUXMIPS64):
	CGO_ENABLED=0 GOOS=linux GOARCH=mips64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUXMIPS64) ./main/
$(LINUXMIPS64LE):
	CGO_ENABLED=0 GOOS=linux GOARCH=mips64le $(GO) build -ldflags $(LDFLAGS) -o $(LINUXMIPS64LE) ./main/
$(LINUXLOONG64):
	CGO_ENABLED=0 GOOS=linux GOARCH=loong64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUXLOONG64) ./main/
$(LINUXRISCV64):
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUXRISCV64) ./main/
$(WIN86).exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GO) build -ldflags $(LDFLAGS) -o $(WIN86).exe ./main/
$(WIN64).exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(WIN64).exe ./main/
$(WINARM).exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm $(GO) build -ldflags $(LDFLAGS) -o $(WINARM).exe ./main/
$(WINARM64).exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GO) build -ldflags $(LDFLAGS) -o $(WINARM64).exe ./main/

.PHONY: cleanzip cleanbin clean upload

cleanzip:
	rm -f qshell-$(VERSION)-*.zip

cleanbin:
	rm -f qshell-$(VERSION)-*

clean: cleanzip cleanbin

upload:
	qshell rput devtools $(DARWIN).zip $(DARWIN).zip
	qshell rput devtools $(DARWINARM64).zip $(DARWINARM64).zip
	qshell rput devtools $(LINUX86).zip $(LINUX86).zip
	qshell rput devtools $(LINUX64).zip $(LINUX64).zip
	qshell rput devtools $(LINUXARM).zip $(LINUXARM).zip
	qshell rput devtools $(LINUXARM64).zip $(LINUXARM64).zip
	qshell rput devtools $(LINUXMIPS).zip $(LINUXMIPS).zip
	qshell rput devtools $(LINUXMIPSLE).zip $(LINUXMIPSLE).zip
	qshell rput devtools $(LINUXMIPS64).zip $(LINUXMIPS64).zip
	qshell rput devtools $(LINUXMIPS64LE).zip $(LINUXMIPS64LE).zip
	qshell rput devtools $(LINUXLOONG64).zip $(LINUXLOONG64).zip
	qshell rput devtools $(LINUXRISCV64).zip $(LINUXRISCV64).zip
	qshell rput devtools $(WIN86).zip $(WIN86).zip
	qshell rput devtools $(WIN64).zip $(WIN64).zip
	qshell rput devtools $(WINARM).zip $(WINARM).zip
	qshell rput devtools $(WINARM64).zip $(WINARM64).zip
