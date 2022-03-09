package utils

import "testing"

func TestIsIp(t *testing.T) {
	ip := "10.200.20.23"
	if !IsIPString(ip) {
		t.Fatal(ip, "should be ip")
	}
}

func TestContainIPInString(t *testing.T) {
	ip := "10.200.20.23:80"
	if !IsIPUrlString(ip) {
		t.Fatal(ip, "should be ip")
	}

	ip = "aaa::bb:80"
	if !IsIPUrlString(ip) {
		t.Fatal(ip, "should be ip")
	}

	ip = "2001:0db8:86a3:08d3:1319:8a2e:0370:7344"
	if !IsIPUrlString(ip) {
		t.Fatal(ip, "should be ip")
	}
}