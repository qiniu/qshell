package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"testing"
)

var accessKey = test.AccessKey
var secretKey = test.SecretKey

func TestUser(t *testing.T) {
	TestUserIntegration(t)
}

func TestUserIntegration(t *testing.T) {
	success, err := cleanUser()
	if len(err) > 0 || !success {
		t.Fatal("clean first not success:", err)
	}

	has, err := hasUser()
	if len(err) > 0 || has {
		t.Fatal("shouldn't has user after clean first:", err)
	}

	// 增加 1
	userName := "test_add_1"
	success, err = addUser(userName, accessKey, secretKey)
	if len(err) > 0 || !success {
		t.Fatal("add user not success:", err)
	}

	has, err = hasUser()
	if len(err) > 0 || !has {
		t.Fatal("should has user after add:", err)
	}

	has, err = containUser(userName)
	if len(err) > 0 || !has {
		t.Fatal("should contain user after add:", err)
	}

	has, err = lookupUser(userName)
	if len(err) > 0 || !has {
		t.Fatal("should lookup user success after add:", err)
	}

	is, err := currentUserIs(userName)
	if len(err) > 0 || !is {
		t.Fatal("should change current user after add:", err)
	}

	// 增加 2
	userName = "test_add_2"
	success, err = addUserWithLongOptions(userName, accessKey, secretKey)
	if len(err) > 0 || !success {
		t.Fatal("add user 2 not success:", err)
	}

	has, err = hasUser()
	if len(err) > 0 || !has {
		t.Fatal("should has user after add 2:", err)
	}

	has, err = containUser(userName)
	if len(err) > 0 || !has {
		t.Fatal("shouldn contain user after add 2:", err)
	}

	has, err = lookupUser(userName)
	if len(err) > 0 || !has {
		t.Fatal("should lookup user success after add 2:", err)
	}

	is, err = currentUserIs(userName)
	if len(err) > 0 || !is {
		t.Fatal("should change current user after add 2:", err)
	}

	// 改
	userName = "test_add_1"
	success, err = changeCurrentUser(userName)
	if len(err) > 0 || !success {
		t.Fatal("change current user not success:", err)
	}

	is, err = currentUserIs(userName)
	if len(err) > 0 || !is {
		t.Fatal("should change current user after change :", err)
	}

	// 删除
	success, err = deleteUser(userName)
	if len(err) > 0 || !success {
		t.Fatal("remove user not success:", err)
	}

	has, err = lookupUser(userName)
	if len(err) > 0 || has {
		t.Fatal("shouldn't lookup user success after remove:", err)
	}

	is, err = currentUserIs(userName)
	if len(err) > 0 || !is {
		t.Fatal("shouldn't change current user after delete current user :", err)
	}

	// 清除
	success, err = cleanUser()
	if len(err) > 0 || !success {
		t.Fatal("clean end not success:", err)
	}

	has, err = hasUser()
	if len(err) > 0 || has {
		t.Fatal("shouldn't has user after clean end:", err)
	}


	// 添加测试账号
	userName = "Kodo"
	success, err = addUser(userName, accessKey, secretKey)
	if len(err) > 0 || !success {
		t.Fatal("add user not success:", err)
	}
}
