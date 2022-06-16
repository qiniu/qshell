//go:build integration

package cmd

import (
	"github.com/qiniu/qshell/v2/cmd_test/test"
	"strings"
)

func addUser(userName, ak, sk string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "add", ak, sk, userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "add error"
	}).Run()
	success = len(result) == 0
	return
}

func addUserWithLongOptions(userName, ak, sk string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "add", "--ak", ak, "--sk", sk, "--name", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "add with long options error"
	}).Run()
	success = len(result) == 0
	return
}

func cleanUser() (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "clean").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "clean user error"
	}).Run()
	success = len(result) == 0
	return
}

func deleteUser(userName string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "remove", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "add user error"
	}).Run()
	success = len(result) == 0
	return
}

func changeCurrentUser(userName string) (success bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "cu", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "change current error"
	}).Run()
	success = strings.Contains(result, "success")
	return
}

func currentUserIs(userName string) (isUser bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "current").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "get current error"
	}).Run()

	isUser = strings.Contains(result, userName)
	return
}

func lookupUser(userName string) (isFound bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "lookup", userName).ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "lookup user error"
	}).Run()
	isFound = strings.Contains(result, userName)
	return
}

func containUser(userName string) (contain bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "ls").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "ls error"
	}).Run()

	contain = strings.Contains(result, userName)
	return
}

func hasUser() (hasUser bool, errorMsg string) {
	result := ""
	test.NewTestFlow("user", "ls").ResultHandler(func(line string) {
		result += line
	}).ErrorHandler(func(line string) {
		errorMsg = "ls error"
	}).Run()
	hasUser = len(result) > 0
	return
}
