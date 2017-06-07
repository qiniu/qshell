package conf

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"

	"qiniu/rpc"
)

var UP_HOST = "http://up.qiniu.com"
var RSF_HOST = "http://rsf.qbox.me"
var IO_HOST = "http://iovip.qbox.me"

var RS_HOST = "http://rs.qiniu.com"
var API_HOST = "http://api.qiniu.com"

var ACCESS_KEY string
var SECRET_KEY string

var version = "6.0.6"

var userPattern = regexp.MustCompile("^[a-zA-Z0-9_.-]*$")

// user should be [A-Za-z0-9]*
func SetUser(user string) error {
	if !userPattern.MatchString(user) {
		return errors.New("invalid user format")
	}
	rpc.UserAgent = formatUserAgent(user)
	return nil
}

func formatUserAgent(user string) string {
	return fmt.Sprintf("QiniuGo/%s (%s; %s; %s) %s", version, runtime.GOOS, runtime.GOARCH, user, runtime.Version())
}

func init() {
	SetUser("")
}
