package version

import (
	"runtime/debug"
)

var (
	defaultVersion = "UNSTABLE"
	version        = defaultVersion
)

func Version() string {
	// 通过 LDFLAGS 设置的版本号
	if version != defaultVersion {
		return version
	}

	v, ok := buildInfoVersion()
	if !ok {
		return version
	}

	
	if v == "(devel)" {
		return version
	}
	return v
}

func buildInfoVersion() (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	return info.Main.Version, true
}
