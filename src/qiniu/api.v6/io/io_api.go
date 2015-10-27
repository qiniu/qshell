package io

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"strconv"
	"strings"

	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
)

// ----------------------------------------------------------

// @gist PutExtra
type PutExtra struct {
	Params map[string]string //可选，用户自定义参数，必须以 "x:" 开头
	//若不以x:开头，则忽略
	MimeType string //可选，当为 "" 时候，服务端自动判断
	Crc32    uint32
	CheckCrc uint32
	// CheckCrc == 0: 表示不进行 crc32 校验
	// CheckCrc == 1: 对于 Put 等同于 CheckCrc = 2；对于 PutFile 会自动计算 crc32 值
	// CheckCrc == 2: 表示进行 crc32 校验，且 crc32 值就是上面的 Crc32 变量
}

// @endgist

type PutRet struct {
	Hash string `json:"hash"` // 如果 uptoken 没有指定 ReturnBody，那么返回值是标准的 PutRet 结构
	Key  string `json:"key"`
}

var tmpFilePrefix = "qiniu-go-sdk-tmpfile"

// ----------------------------------------------------------
// !!! Deprecated !!!
//
// 1. 不推荐使用该组 API, 因为可能造成本地磁盘IO
// 2. 如果只是纯粹上传一个 io.Reader, 请使用 Put2 或者 PutWithoutKey2
// 3. 如果需要上传一个文件, 请使用 PutFile 或者 PutFileWithoutKey

func Put(l rpc.Logger, ret interface{}, uptoken, key string, data io.Reader, extra *PutExtra) error {
	return putReader(l, ret, uptoken, key, true, data, extra)
}

func PutWithoutKey(l rpc.Logger, ret interface{}, uptoken string, data io.Reader, extra *PutExtra) error {
	return putReader(l, ret, uptoken, "", false, data, extra)
}

func putReader(l rpc.Logger, ret interface{}, uptoken, key string, hasKey bool, data io.Reader, extra *PutExtra) error {

	rs, ok := data.(io.ReadSeeker)
	if ok {
		// 通过 Seeker 接口获取大小
		size, err := rs.Seek(0, 2)
		if err != nil {
			return err
		}
		_, err = rs.Seek(0, 0)
		if err != nil {
			return err
		}
		return put(l, ret, uptoken, key, hasKey, data, size, extra)
	} else {
		// 写临时文件
		tmpf, err := ioutil.TempFile(os.TempDir(), tmpFilePrefix)
		if err != nil {
			return err
		}
		fname := tmpf.Name()
		defer os.Remove(fname)

		_, err = io.Copy(tmpf, data)
		if err != nil {
			tmpf.Close()
			return err
		}
		tmpf.Close()

		return putFile(l, ret, uptoken, key, hasKey, fname, extra)
	}

	return nil
}

// ----------------------------------------------------------

func put(l rpc.Logger, ret interface{}, uptoken, key string, hasKey bool, data io.Reader, size int64, extra *PutExtra) error {

	// CheckCrc == 1: 对于 Put 和 PutWithoutKey 等同于 CheckCrc == 2
	if extra != nil {
		if extra.CheckCrc == 1 {
			extra1 := *extra
			extra = &extra1
			extra.CheckCrc = 2
		}
	}
	return putWrite(l, ret, uptoken, key, hasKey, data, size, extra)
}

func Put2(l rpc.Logger, ret interface{}, uptoken, key string, data io.Reader, size int64, extra *PutExtra) error {
	return put(l, ret, uptoken, key, true, data, size, extra)
}

func PutWithoutKey2(l rpc.Logger, ret interface{}, uptoken string, data io.Reader, size int64, extra *PutExtra) error {
	return put(l, ret, uptoken, "", false, data, size, extra)
}

// ----------------------------------------------------------

func PutFile(l rpc.Logger, ret interface{}, uptoken, key, localFile string, extra *PutExtra) (err error) {
	return putFile(l, ret, uptoken, key, true, localFile, extra)
}

func PutFileWithoutKey(l rpc.Logger, ret interface{}, uptoken, localFile string, extra *PutExtra) (err error) {
	return putFile(l, ret, uptoken, "", false, localFile, extra)
}

func putFile(l rpc.Logger, ret interface{}, uptoken, key string, hasKey bool, localFile string, extra *PutExtra) (err error) {

	f, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return
	}
	fsize := finfo.Size()

	if extra != nil && extra.CheckCrc == 1 {
		extra.Crc32, err = getFileCrc32(f)
		if err != nil {
			return
		}
	}

	return putWrite(l, ret, uptoken, key, hasKey, f, fsize, extra)
}

// ----------------------------------------------------------

func putWrite(l rpc.Logger, ret interface{}, uptoken, key string, hasKey bool, data io.Reader, size int64, extra *PutExtra) error {

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	err := writeMultipart(writer, uptoken, key, hasKey, extra)
	if err != nil {
		return err
	}

	lastLine := fmt.Sprintf("\r\n--%s--\r\n", writer.Boundary())
	r := bytes.NewReader([]byte(lastLine))

	bodyLen := int64(b.Len()) + size + int64(len(lastLine))
	mr := io.MultiReader(&b, data, r)

	contentType := writer.FormDataContentType()
	return rpc.DefaultClient.CallWith64(l, ret, UP_HOST, contentType, mr, bodyLen)
}

/*
 * extra.CheckCrc:
 *      0:     不进行crc32校验
 *      1:     以writeMultipart自动生成crc32的值，进行校验
 *      2:     以extra.Crc32的值，进行校验
 *      other: 和2一样， 以 extra.Crc32的值，进行校验
 */
func writeMultipart(writer *multipart.Writer, uptoken, key string, hasKey bool, extra *PutExtra) (err error) {

	if extra == nil {
		extra = &PutExtra{}
	}

	//token
	if err = writer.WriteField("token", uptoken); err != nil {
		return
	}

	//key
	if hasKey {
		if err = writer.WriteField("key", key); err != nil {
			return
		}
	}

	// extra.Params
	if extra.Params != nil {
		for k, v := range extra.Params {
			err = writer.WriteField(k, v)
			if err != nil {
				return
			}
		}
	}

	//extra.CheckCrc
	if extra.CheckCrc != 0 {
		err = writer.WriteField("crc32", strconv.FormatInt(int64(extra.Crc32), 10))
		if err != nil {
			return
		}
	}

	//file
	head := make(textproto.MIMEHeader)

	// default the filename is same as key , but ""
	var fileName = key
	if fileName == "" {
		fileName = "filename"
	}

	head.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(fileName)))
	if extra.MimeType != "" {
		head.Set("Content-Type", extra.MimeType)
	}

	_, err = writer.CreatePart(head)
	return err
}

// ----------------------------------------------------------

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func getFileCrc32(f *os.File) (uint32, error) {
	defer f.Seek(0, 0)

	h := crc32.NewIEEE()
	_, err := io.Copy(h, f)

	return h.Sum32(), err
}
