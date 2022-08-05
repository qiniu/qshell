//go:build integration

package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"testing"
)

func TestCmd(t *testing.T) {
	data.SetTestMode()

	TestUser(t)
	ClearCache(t)
}

func ClearCache(t *testing.T) {
	path := test.RemoveRootPath()
	if err := path; err != nil {
		fmt.Printf("Remove Cache Path:%s error:%v", path, err)
	}

	path = test.RemoveTempPath()
	if err := path; err != nil {
		fmt.Printf("Remove Cache Path:%s error:%v", path, err)
	}

	path = test.RemoveResultPath()
	if err := path; err != nil {
		fmt.Printf("Remove Cache Path:%s error:%v", path, err)
	}
}
