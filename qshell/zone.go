package qshell

type ZoneConfig struct {
	UpHost    string
	RsHost    string
	RsfHost   string
	IovipHost string
	ApiHost   string
}

const (
	ZoneNB  = "z0"
	ZoneBC  = "z1"
	ZoneHN  = "z2"
	ZoneNA0 = "na0"
	ZoneAS0 = "as0"
)

//zone all defaults to the service source site

var ZoneNBConfig = ZoneConfig{
	UpHost:    "http://upload.qiniu.com",
	IovipHost: "http://iovip.qbox.me",
	RsHost:    "http://rs.qiniu.com",
	RsfHost:   "http://rsf.qiniu.com",
	ApiHost:   "http://api.qiniu.com",
}

var ZoneBCConfig = ZoneConfig{
	UpHost:    "http://upload-z1.qiniu.com",
	IovipHost: "http://iovip-z1.qbox.me",
	RsHost:    "http://rs-z1.qiniu.com",
	RsfHost:   "http://rsf-z1.qiniu.com",
	ApiHost:   "http://api-z1.qiniu.com",
}

var ZoneHNConfig = ZoneConfig{
	UpHost:    "http://upload-z2.qiniu.com",
	IovipHost: "http://iovip-z2.qbox.me",
	RsHost:    "http://rs-z2.qiniu.com",
	RsfHost:   "http://rsf-z2.qiniu.com",
	ApiHost:   "http://api-z2.qiniu.com",
}

var ZoneNA0Config = ZoneConfig{
	UpHost:    "http://upload-na0.qiniu.com",
	IovipHost: "http://iovip-na0.qbox.me",
	RsHost:    "http://rs-na0.qiniu.com",
	RsfHost:   "http://rsf-na0.qiniu.com",
	ApiHost:   "http://api-na0.qiniu.com",
}

var ZoneSA0Config = ZoneConfig{
	UpHost:    "http://upload-as0.qiniu.com",
	IovipHost: "http://iovip-as0.qbox.me",
	RsHost:    "http://rs-as0.qiniu.com",
	RsfHost:   "http://rsf-as0.qiniu.com",
	ApiHost:   "http://api-as0.qiniu.com",
}

func SetZone(zone string) {
	var zoneConfig ZoneConfig
	switch zone {
	case ZoneBC:
		zoneConfig = ZoneBCConfig
	case ZoneHN:
		zoneConfig = ZoneHNConfig
	case ZoneNA0:
		zoneConfig = ZoneNA0Config
	case ZoneAS0:
		zoneConfig = ZoneSA0Config
	default:
		zoneConfig = ZoneNBConfig
	}
	if hosts[UpHostId] == "" {
		hosts[UpHostId] = zoneConfig.UpHost
	}
	if hosts[RsHostId] == "" {
		hosts[RsHostId] = zoneConfig.RsHost
	}
	if hosts[RsfHostId] == "" {
		hosts[RsfHostId] = zoneConfig.RsfHost
	}
	if hosts[IoHostId] == "" {
		hosts[IoHostId] = zoneConfig.IovipHost
	}
	if hosts[ApiHostId] == "" {
		hosts[ApiHostId] = zoneConfig.ApiHost
	}
}

func IsValidZone(zone string) (valid bool) {
	switch zone {
	case ZoneNB, ZoneBC, ZoneHN, ZoneNA0, ZoneAS0:
		valid = true
	default:
		valid = false
	}
	return
}
