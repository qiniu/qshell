package object

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"net/url"
)

type SaveAsApiInfo struct {
	PublicUrl  string `json:"public_url"`
	SaveBucket string `json:"save_bucket"`
	SaveKey    string `json:"save_key"`
}

func SaveAs(info SaveAsApiInfo) (string, *data.CodeError) {
	if len(info.PublicUrl) == 0 {
		return "", alert.CannotEmptyError("public url", "")
	}

	uri, parseErr := url.Parse(info.PublicUrl)
	if parseErr != nil {
		return "", data.NewEmptyError().AppendDesc("parse public url error:" + parseErr.Error())
	}

	if len(info.SaveBucket) == 0 {
		return "", alert.CannotEmptyError("save bucket", "")
	}

	if len(info.SaveBucket) == 0 {
		return "", alert.CannotEmptyError("save key", "")
	}

	acc, err := workspace.GetAccount()
	if err != nil {
		return "", err
	}

	baseUrl := uri.Host + uri.RequestURI()
	saveEntry := info.SaveBucket + ":" + info.SaveKey
	encodedSaveEntry := base64.URLEncoding.EncodeToString([]byte(saveEntry))
	baseUrl += "|saveas/" + encodedSaveEntry
	h := hmac.New(sha1.New, []byte(acc.SecretKey))
	h.Write([]byte(baseUrl))
	sign := h.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	return info.PublicUrl + "|saveas/" + encodedSaveEntry + "/sign/" + acc.AccessKey + ":" + encodedSign, nil
}
