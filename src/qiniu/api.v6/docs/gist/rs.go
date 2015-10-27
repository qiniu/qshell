package gist

import (
	"log"

	"github.com/qiniu/api/rs"

	. "github.com/qiniu/api/conf"
)

func init() {
	// @gist init
	ACCESS_KEY = "<YOUR_APP_ACCESS_KEY>"
	SECRET_KEY = "<YOUR_APP_SECRET_KEY>"
	// @endgist
}

func rsDemo(bucket, key, bucketSrc, keySrc, bucketDest, keyDest string) {

	// @gist rsPre
	//此操作前 请确保 accesskey和secretkey 已被正确赋值
	var rsCli = rs.New(nil)
	var err error
	// @endgist

	// @gist rsStat
	var ret rs.Entry
	ret, err = rsCli.Stat(nil, bucket, key)
	if err != nil {
		// 产生错误
		log.Println("rs.Stat failed:", err)
		return
	}
	// 处理返回值
	log.Println(ret)
	// @endgist

	// @gist rsCopy
	err = rsCli.Copy(nil, bucketSrc, keySrc, bucketDest, keyDest)
	if err != nil {
		// 产生错误
		log.Println("rs.Copy failed:", err)
		return
	}
	// @endgist

	// @gist rsMove
	err = rsCli.Move(nil, bucketSrc, keySrc, bucketDest, keyDest)
	if err != nil {
		//产生错误
		log.Println("rs.Copy failed:", err)
		return
	}
	// @endgist

	// @gist rsDelete
	err = rsCli.Delete(nil, bucket, key)
	if err != nil {
		// 产生错误
		log.Println("rs.Copy failed:", err)
		return
	}
	// @endgist
}

func batchDemo(bucket, key, bucket1, key1, bucket2, key2, bucket3, key3, bucket4, key4 string) {

	// @gist rsBatchPre
	// 此操作前 请确保 accesskey和secretkey 已被正确赋值
	var rsCli = rs.New(nil)
	var err error
	// @endgist

	// @gist rsEntryPathes
	entryPathes := []rs.EntryPath{
		rs.EntryPath{
			Bucket: bucket1,
			Key:    key1,
		},
		rs.EntryPath{
			Bucket: bucket2,
			Key:    key2,
		},
	}
	// @endgist

	// @gist rsPathPairs
	// 每个复制操作都含有源文件和目标文件
	entryPairs := []rs.EntryPathPair{
		rs.EntryPathPair{
			Src: rs.EntryPath{
				Bucket: bucket1,
				Key:    key1,
			},
			Dest: rs.EntryPath{
				Bucket: bucket2,
				Key:    key2,
			},
		}, rs.EntryPathPair{
			Src: rs.EntryPath{
				Bucket: bucket3,
				Key:    key3,
			},
			Dest: rs.EntryPath{
				Bucket: bucket4,
				Key:    key4,
			},
		},
	}
	// @endgist

	// @gist rsBatchStat
	var batchStatRets []rs.BatchStatItemRet
	batchStatRets, err = rsCli.BatchStat(nil, entryPathes) // []rs.BatchStatItemRet, error
	if err != nil {
		// 产生错误
		log.Println("rs.BatchStat failed:", err)
		return
	}
	// 处理返回值
	for _, item := range batchStatRets {
		log.Println(item)
	}
	// @endgist

	// @gist rsBatchCopy
	var batchCopyRets []rs.BatchItemRet
	batchCopyRets, err = rsCli.BatchCopy(nil, entryPairs)
	if err != nil {
		// 产生错误
		log.Println("rs.BatchCopy failed:", err)
		return
	}
	for _, item := range batchCopyRets {
		// 遍历每个操作的返回结果
		log.Println(item.Code, item.Error)
	}
	// @endgist

	// @gist rsBatchMove
	var batchMoveRets []rs.BatchItemRet
	batchMoveRets, err = rsCli.BatchMove(nil, entryPairs)
	if err != nil {
		// 产生错误
		log.Println("rs.BatchMove failed:", err)
		return
	}
	for _, item := range batchMoveRets {
		// 遍历每个操作的返回结果
		log.Println(item.Code, item.Error)
	}
	// @endgist

	// @gist rsBatchDelete
	var batchDeleteRets []rs.BatchItemRet
	batchDeleteRets, err = rsCli.BatchDelete(nil, entryPathes)
	if err != nil {
		// 产生错误
		log.Println("rs.BatchDelete failed:", err)
		return
	}
	for _, item := range batchDeleteRets {
		// 遍历每个操作的返回结果
		log.Println(item.Code, item.Error)
	}
	// @endgist

	// @gist rsBatchAdv
	ops := []string{
		rs.URIStat(bucket, key1),
		rs.URICopy(bucket, key1, bucket, key2), // 复制key1到key2
		rs.URIDelete(bucket, key1),             // 删除key1
		rs.URIMove(bucket, key2, bucket, key1), //将key2移动到key1
	}

	rets := new([]rs.BatchItemRet)
	err = rsCli.Batch(nil, rets, ops)
	if err != nil {
		// 产生错误
		log.Println("rs.Batch failed:", err)
		return
	}
	for _, ret := range *rets {
		log.Println(ret.Code, ret.Error)
	}
	// @endgist
}
