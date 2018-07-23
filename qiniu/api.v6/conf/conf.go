package conf

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"

	"github.com/tonycai653/iqshell/qiniu/rpc"
)

var (
	UP_HOST         = "http://up.qiniu.com"
	RSF_HOST        = "http://rsf.qbox.me"
	IO_HOST         = "http://iovip.qbox.me"
	RS_HOST         = "http://rs.qiniu.com"
	API_HOST        = "http://api.qiniu.com"
	BUCKET_RS_HOST  = "http://rs.qiniu.com"
	BUCKET_API_HOST = "http://api.qiniu.com"
)

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
