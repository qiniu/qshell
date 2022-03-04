package utils

import "testing"

func TestIsCmdExist(t *testing.T) {
	exist := IsCmdExist("ls")
	if !exist {
		t.Fatal("ls should exist")
	}

	exist = IsCmdExist("lss")
	if exist {
		t.Fatal("lls shouldn't exist")
	}
}
