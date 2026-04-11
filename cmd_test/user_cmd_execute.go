//go:build integration

package cmd

import (
	"strings"

	"github.com/qiniu/qshell/v2/cmd_test/test"
)

func addUser(userName, ak, sk string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "add", ak, sk, userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "add error"
	}).Run()
	success = len(result) == 0
	return success, errorMsg
}

func addUserWithLongOptions(userName, ak, sk string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "add", "--ak", ak, "--sk", sk, "--name", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "add with long options error"
	}).Run()
	success = len(result) == 0
	return success, errorMsg
}

func cleanUser() (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "clean").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "clean user error"
	}).Run()
	success = len(result) == 0
	return success, errorMsg
}

func deleteUser(userName string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "remove", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "add user error"
	}).Run()
	success = len(result) == 0
	return success, errorMsg
}

func changeCurrentUser(userName string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "cu", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "change current error"
	}).Run()
	success = strings.Contains(result, "success")
	return success, errorMsg
}

func currentUserIs(userName string) (isUser bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "current").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "get current error"
	}).Run()

	isUser = strings.Contains(result, userName)
	return isUser, errorMsg
}

func lookupUser(userName string) (isFound bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "lookup", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "lookup user error"
	}).Run()
	isFound = strings.Contains(result, userName)
	return isFound, errorMsg
}

func containUser(userName string) (contain bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "ls").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "ls error"
	}).Run()

	contain = strings.Contains(result, userName)
	return contain, errorMsg
}

func hasUser() (hasUser bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "ls").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "ls error"
	}).Run()
	hasUser = len(result) > 0
	return hasUser, errorMsg
}
