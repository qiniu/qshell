package qshell

import (
	"qiniu/api.v6/conf"
)

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

func SetZone(zone string) {
	var zoneConfig ZoneConfig
	switch zone {
	case ZoneBC:
		zoneConfig = ZoneBCConfig
	case ZoneHN:
		zoneConfig = ZoneHNConfig
	case ZoneNA0:
		zoneConfig = ZoneNA0Config
	default:
		zoneConfig = ZoneNBConfig
	}
	conf.UP_HOST = zoneConfig.UpHost
	conf.RS_HOST = zoneConfig.RsHost
	conf.RSF_HOST = zoneConfig.RsfHost
	conf.IO_HOST = zoneConfig.IovipHost
	conf.API_HOST = zoneConfig.ApiHost
}

func IsValidZone(zone string) (valid bool) {
	switch zone {
	case ZoneNB, ZoneBC, ZoneHN, ZoneNA0:
		valid = true
	default:
		valid = false
	}
	return
}
