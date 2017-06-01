package rs

import (
	"testing"
)

func init() {

	client = New(nil)
	// 删除 可能存在的 newkey1  newkey2
	client.Delete(nil, bucketName, key)
	client.Delete(nil, bucketName, newkey1)
	client.Delete(nil, bucketName, newkey2)
}

func TestAll(t *testing.T) {

	//上传一个文件用用于测试
	err := upFile("batch_api.go", bucketName, key)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Delete(nil, bucketName, key)

	testBatchStat(t)
	testBatchCopy(t)
	testBatchMove(t)
	testBatchDelete(t)
	testBatch(t)
}

func testBatchStat(t *testing.T) {

	entryPath := EntryPath{
		Bucket: bucketName,
		Key:    key,
	}

	rets, err := client.BatchStat(nil, []EntryPath{entryPath, entryPath, entryPath})
	if err != nil {
		t.Fatal(err)
	}

	if len(rets) != 3 {
		t.Fatal("BatchStat failed: len(result) = ", len(rets))
	}

	stat, _ := client.Stat(nil, bucketName, key)

	if rets[0].Data != stat || rets[1].Data != stat || rets[2].Data != stat {
		t.Fatal("BatchStat failed : returns err")
	}
}

func testBatchMove(t *testing.T) {

	stat0, err := client.Stat(nil, bucketName, key)
	if err != nil {
		t.Fatal("BathMove get stat failed:", err)
	}
	entryPair1 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    key,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
	}

	entryPair2 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey2,
		},
	}

	_, err = client.BatchMove(nil, []EntryPathPair{entryPair1, entryPair2})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Move(nil, bucketName, newkey2, bucketName, key)

	stat1, err := client.Stat(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal("BathMove get stat failed:", err)
	}

	if stat0.Hash != stat1.Hash {
		t.Fatal("BatchMove failed : Move err", stat0, stat1)
	}
}

func testBatchCopy(t *testing.T) {

	entryPair1 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    key,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
	}

	entryPair2 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey2,
		},
	}

	_, err := client.BatchCopy(nil, []EntryPathPair{entryPair1, entryPair2})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Delete(nil, bucketName, newkey1)
	defer client.Delete(nil, bucketName, newkey2)

	stat0, _ := client.Stat(nil, bucketName, key)
	stat1, _ := client.Stat(nil, bucketName, newkey1)
	stat2, _ := client.Stat(nil, bucketName, newkey2)
	if stat0.Hash != stat1.Hash || stat0.Hash != stat2.Hash {
		t.Fatal("BatchCopy failed : Copy err")
	}
}

func testBatchDelete(t *testing.T) {

	client.Copy(nil, bucketName, key, bucketName, newkey1)
	client.Copy(nil, bucketName, key, bucketName, newkey2)

	entryPath1 := EntryPath{
		Bucket: bucketName,
		Key:    newkey1,
	}
	entryPath2 := EntryPath{
		Bucket: bucketName,
		Key:    newkey2,
	}

	_, err := client.BatchDelete(nil, []EntryPath{entryPath1, entryPath2})
	if err != nil {
		t.Fatal(err)
	}

	_, err1 := client.Stat(nil, bucketName, newkey1)
	_, err2 := client.Stat(nil, bucketName, newkey2)

	//这里 err1 != nil，否则文件没被成功删除
	if err1 == nil || err2 == nil {
		t.Fatal("BatchDelete failed : File do not delete")
	}
}

func testBatch(t *testing.T) {

	ops := []string{
		URICopy(bucketName, key, bucketName, newkey1),
		URIDelete(bucketName, key),
		URIMove(bucketName, newkey1, bucketName, key),
	}
	rets := new([]BatchItemRet)
	err := client.Batch(nil, rets, ops)
	if err != nil {
		t.Fatal(err)
	}
}
