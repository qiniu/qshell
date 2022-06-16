package utils

import (
	"testing"
)

func TestBytesToReadable(t *testing.T) {
	sizes := map[int64]string{
		512:           "512B",
		1024:          "1.00KB",
		2048:          "2.00KB",
		1048576:       "1.00MB",
		1073741824:    "1.00GB",
		2073741824:    "1.93GB",
		1099511627776: "1.00TB",
	}

	for size, want := range sizes {
		got := BytesToReadable(size)
		if got != want {
			t.Fatalf("size got=%s, want=%s", got, want)
		}
	}
}

func TestKeyFromUrl(t *testing.T) {
	url := "http://vod4a6mk39q.nosdn.127.net/b258912a66334476851b698d6fe64931_1558331445602_1558331488089_2062207192-00000.mp4?download=%E7%A7%80%E7%9B%B4%E6%92%AD%E7%BC%96%E5%8F%B72114_20190520-135045_20190520-135128.mp4"
	want := "b258912a66334476851b698d6fe64931_1558331445602_1558331488089_2062207192-00000.mp4"

	key, err := KeyFromUrl(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if key != want {
		t.Fatalf("got = %s, want = %s\n", key, want)
	}
}

func TestRemoveUrlScheme(t *testing.T) {
	host := "hqiniu.com"
	url := host
	result := RemoveUrlScheme(url)
	if host != result {
		t.Fatalf("RemoveUrlScheme failed, excpet:%s but:%s\n", host, result)
	}

	url = "http://" + host
	result = RemoveUrlScheme(url)
	if host != result {
		t.Fatalf("RemoveUrlScheme http:// failed, excpet:%s but:%s\n", host, result)
	}

	url = "https://" + host
	result = RemoveUrlScheme(url)
	if host != result {
		t.Fatalf("RemoveUrlScheme https:// failed, excpet:%s but:%s\n", host, result)
	}
}
