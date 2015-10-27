package rs

import (
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
)

// ----------------------------------------------------------

func (rs Client) Batch(l rpc.Logger, ret interface{}, op []string) (err error) {
	return rs.Conn.CallWithForm(l, ret, RS_HOST+"/batch", map[string][]string{"op": op})
}

// ----------------------------------------------------------

// @gist batchStatItemRet
type BatchStatItemRet struct {
	Data  Entry  `json:"data"`
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// @endgist

// @gist entryPath
type EntryPath struct {
	Bucket string
	Key    string
}

// @endgist

func (rs Client) BatchStat(l rpc.Logger, entries []EntryPath) (ret []BatchStatItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIStat(e.Bucket, e.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

// ----------------------------------------------------------

// @gist batchItemRet
type BatchItemRet struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// @endgist

func (rs Client) BatchDelete(l rpc.Logger, entries []EntryPath) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIDelete(e.Bucket, e.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

// ----------------------------------------------------------

// @gist entryPathPair
type EntryPathPair struct {
	Src  EntryPath
	Dest EntryPath
}

// @endgist

func (rs Client) BatchMove(l rpc.Logger, entries []EntryPathPair) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIMove(e.Src.Bucket, e.Src.Key, e.Dest.Bucket, e.Dest.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

func (rs Client) BatchCopy(l rpc.Logger, entries []EntryPathPair) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URICopy(e.Src.Bucket, e.Src.Key, e.Dest.Bucket, e.Dest.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

// ----------------------------------------------------------
