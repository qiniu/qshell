package cli

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/log"
	"net/url"
	"qshell"
	"strconv"
	"time"
)

func Base64Encode(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		urlSafe := true
		var err error
		var dataToEncode string
		if len(params) == 2 {
			urlSafe, err = strconv.ParseBool(params[0])
			if err != nil {
				log.Error("Invalid bool value or <UrlSafe>, must true or false")
				return
			}
			dataToEncode = params[1]
		} else {
			dataToEncode = params[0]
		}
		dataEncoded := ""
		if urlSafe {
			dataEncoded = base64.URLEncoding.EncodeToString([]byte(dataToEncode))
		} else {
			dataEncoded = base64.StdEncoding.EncodeToString([]byte(dataToEncode))
		}
		fmt.Println(dataEncoded)
	} else {
		CmdHelp(cmd)
	}
}
func Base64Decode(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		urlSafe := true
		var err error
		var dataToDecode string
		if len(params) == 2 {
			urlSafe, err = strconv.ParseBool(params[0])
			if err != nil {
				log.Error("Invalid bool value or <UrlSafe>, must true or false")
				return
			}
			dataToDecode = params[1]
		} else {
			dataToDecode = params[0]
		}
		var dataDecoded []byte
		if urlSafe {
			dataDecoded, err = base64.URLEncoding.DecodeString(dataToDecode)
			if err != nil {
				log.Error("Failed to decode `", dataToDecode, "' in url safe mode.")
				return
			}
		} else {
			dataDecoded, err = base64.StdEncoding.DecodeString(dataToDecode)
			if err != nil {
				log.Error("Failed to decode `", dataToDecode, "' in standard mode.")
				return
			}
		}
		fmt.Println(string(dataDecoded))
	} else {
		CmdHelp(cmd)
	}
}

func Timestamp2Date(cmd string, params ...string) {
	if len(params) == 1 {
		ts, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			log.Error("Invalid timestamp value,", params[0])
			return
		}
		t := time.Unix(ts, 0)
		fmt.Println(t.String())
	} else {
		CmdHelp(cmd)
	}
}

func TimestampNano2Date(cmd string, params ...string) {
	if len(params) == 1 {
		tns, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			log.Error("Invalid nano timestamp value,", params[0])
			return
		}
		t := time.Unix(0, tns*100)
		fmt.Println(t.String())
	} else {
		CmdHelp(cmd)
	}
}

func TimestampMilli2Date(cmd string, params ...string) {
	if len(params) == 1 {
		tms, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			log.Error("Invalid milli timestamp value,", params[0])
			return
		}
		t := time.Unix(tms/1000, 0)
		fmt.Println(t.String())
	} else {
		CmdHelp(cmd)
	}
}

func Date2Timestamp(cmd string, params ...string) {
	if len(params) == 1 {
		secs, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			log.Error("Invalid seconds to now,", params[0])
			return
		}
		t := time.Now()
		t = t.Add(time.Second * time.Duration(secs))
		fmt.Println(t.Unix())
	} else {
		CmdHelp(cmd)
	}
}

func Urlencode(cmd string, params ...string) {
	if len(params) == 1 {
		dataToEncode := params[0]
		dataEncoded := url.QueryEscape(dataToEncode)
		fmt.Println(dataEncoded)
	} else {
		CmdHelp(cmd)
	}
}

func Urldecode(cmd string, params ...string) {
	if len(params) == 1 {
		dataToDecode := params[0]
		dataDecoded, err := url.QueryUnescape(dataToDecode)
		if err != nil {
			log.Error("Failed to unescape data `", dataToDecode, "'")
		} else {
			fmt.Println(dataDecoded)
		}
	} else {
		CmdHelp(cmd)
	}
}

func Qetag(cmd string, params ...string) {
	if len(params) == 1 {
		localFilePath := params[0]
		qetag, err := qshell.GetEtag(localFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(qetag)
	} else {
		CmdHelp(cmd)
	}
}
