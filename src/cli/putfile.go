package cli

import (
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/log"
)

func FormPut(cmd string, params ...string) {
	if len(params) == 3 || len(params) == 4 {
		bucket := params[0]
		key := params[1]
		localFile := params[2]
		mimeType := ""
		if len(params) == 4 {
			mimeType = params[3]
		}
		accountS.Get()
		mac := digest.Mac{accountS.AccessKey, []byte(accountS.SecretKey)}
		policy := rs.PutPolicy{}
		policy.Scope = bucket
		putExtra := io.PutExtra{}
		if mimeType != "" {
			putExtra.MimeType = mimeType
		}
		uptoken := policy.Token(&mac)
		putRet := io.PutRet{}
		err := io.PutFile(nil, &putRet, uptoken, key, localFile, &putExtra)
		if err != nil {
			log.Error("Put file error", err)
		} else {
			fmt.Println("Put file", localFile, "=>", bucket, ":", putRet.Key, "(", putRet.Hash, ")", "success!")
		}
	} else {
		Help(cmd)
	}
}

func ResumablePut(cmd string, params ...string) {

}
