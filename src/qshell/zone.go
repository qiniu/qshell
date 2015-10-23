package qshell

type ZoneConfig struct {
	UpHost    string
	RsHost    string
	RsfHost   string
	IovipHost string
	ApiHost   string
}

var (
	DEFAULT_API_HOST = ZoneNBConfig.ApiHost
)

const (
	ZoneNB  = "nb"
	ZoneBC  = "bc"
	ZoneAWS = "aws"
)

var ZoneNBConfig = ZoneConfig{
	UpHost:    "http://up.qiniu.com",
	RsHost:    "http://rs.qiniu.com",
	RsfHost:   "http://rsf.qbox.me",
	IovipHost: "http://iovip.qbox.me",
	ApiHost:   "http://api.qiniu.com",
}

var ZoneBCConfig = ZoneConfig{
	UpHost:    "http://up-z1.qiniu.com",
	RsHost:    "http://rs-z1.qiniu.com",
	RsfHost:   "http://rsf-z1.qbox.me",
	IovipHost: "http://iovip-z1.qbox.me",
	ApiHost:   "http://api-z1.qiniu.com",
}

var ZoneAWSConfig = ZoneConfig{
	UpHost:    "http://up.gdipper.com",
	RsHost:    "http://rs.gdipper.com",
	RsfHost:   "http://rsf.gdipper.com",
	IovipHost: "http://iovip.gdipper.me",
	ApiHost:   "http://api.gdipper.com",
}
