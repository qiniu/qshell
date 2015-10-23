package qshell

import (
	"encoding/json"
	"fmt"
	"github.com/qiniu/rpc"
)

type RSFop struct {
	Account
}

type FopRet struct {
	Id          string `json:"id"`
	Code        int    `json:"code"`
	Desc        string `json:"desc"`
	InputBucket string `json:"inputBucket,omitempty"`
	InputKey    string `json:"inputKey,omitempty"`
	Pipeline    string `json:"pipeline,omitempty"`
	Reqid       string `json:"reqid,omitempty"`
	Items       []FopResult
}

func (this *FopRet) String() string {
	strData := fmt.Sprintf("Id: %s\r\nCode: %d\r\nDesc: %s\r\n", this.Id, this.Code, this.Desc)
	if this.InputBucket != "" {
		strData += fmt.Sprintln(fmt.Sprintf("InputBucket: %s", this.InputBucket))
	}
	if this.InputKey != "" {
		strData += fmt.Sprintln(fmt.Sprintf("InputKey: %s", this.InputKey))
	}
	if this.Pipeline != "" {
		strData += fmt.Sprintln(fmt.Sprintf("Pipeline: %s", this.Pipeline))
	}
	if this.Reqid != "" {
		strData += fmt.Sprintln(fmt.Sprintf("Reqid: %s", this.Reqid))
	}

	strData = fmt.Sprintln(strData)
	for _, item := range this.Items {
		strData += fmt.Sprintf("\tCmd:\t%s\r\n\tCode:\t%d\r\n\tDesc:\t%s\r\n", item.Cmd, item.Code, item.Desc)
		if item.Error != "" {
			strData += fmt.Sprintf("\tError:\t%s\r\n", item.Error)
		} else {
			if item.Hash != "" {
				strData += fmt.Sprintf("\tHash:\t%s\r\n", item.Hash)
			}
			if item.Key != "" {
				strData += fmt.Sprintf("\tKey:\t%s\r\n", item.Key)
			}
			if item.Keys != nil {
				if len(item.Keys) > 0 {
					strData += "\tKeys: {\r\n"
					for _, key := range item.Keys {
						strData += fmt.Sprintf("\t\t%s\r\n", key)
					}
					strData += "\t}\r\n"
				}
			}
		}
		strData += "\r\n"
	}
	return strData
}

type FopResult struct {
	Cmd   string   `json:"cmd"`
	Code  int      `json:"code"`
	Desc  string   `json:"desc"`
	Error string   `json:"error,omitempty"`
	Hash  string   `json:"hash,omitempty"`
	Key   string   `json:"key,omitempty"`
	Keys  []string `json:"keys,omitempty"`
}

func (this *RSFop) Prefop(persistentId string, fopRet *FopRet) (err error) {
	client := rpc.DefaultClient
	resp, respErr := client.Get(nil, fmt.Sprintf("http://%s/status/get/prefop?id=%s", DEFAULT_API_HOST, persistentId))
	if respErr != nil {
		err = respErr
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 == 2 {
		if fopRet != nil && resp.ContentLength != 0 {
			pErr := json.NewDecoder(resp.Body).Decode(fopRet)
			if pErr != nil {
				err = pErr
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return rpc.ResponseError(resp)
}
