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
	ZoneNB  = "nb"
	ZoneBC  = "bc"
	ZoneHN  = "hn"
	ZoneAWS = "aws"
	ZoneNA0 = "na0"
)

//zone all defaults to the service source site

var ZoneNBConfig = ZoneConfig{
	UpHost:    "http://up.qiniu.com",
	RsHost:    "http://rs.qiniu.com",
	RsfHost:   "http://rsf.qbox.me",
	IovipHost: "http://iovip.qbox.me",
	ApiHost:   "http://api.qiniu.com",
}

var ZoneBCConfig = ZoneConfig{
	UpHost:    "http://up-z1.qiniu.com",
	RsHost:    "http://rs.qiniu.com",
	RsfHost:   "http://rsf-z1.qbox.me",
	IovipHost: "http://iovip-z1.qbox.me",
	ApiHost:   "http://api.qiniu.com",
}

var ZoneHNConfig = ZoneConfig{
	UpHost:    "http://up-z2.qiniu.com",
	RsHost:    "http://rs-z2.qiniu.com",
	RsfHost:   "http://rsf-z2.qbox.me",
	IovipHost: "http://iovip-z2.qbox.me",
	ApiHost:   "http://api.qiniu.com",
}

var ZoneNA0Config = ZoneConfig{
	UpHost:    "http://up-na0.qiniu.com",
	RsHost:    "http://rs-na0.qbox.me",
	RsfHost:   "http://rsf-na0.qbox.me",
	IovipHost: "http://iovip-na0.qbox.me",
	ApiHost:   "http://api.qiniu.com",
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
