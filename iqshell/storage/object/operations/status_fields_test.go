package operations

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBatchStatusFields 验证默认输出、字段顺序、大小写匹配及非法字段校验。
func TestBatchStatusFields(t *testing.T) {
	for _, tc := range []struct {
		name   string
		fields string
		want   string
	}{
		{"default", "", "Key,Fsize,Hash,MimeType,PutTime,Type"},
		{"ordered", "expiration,key,MD5", "Expiration,Key,MD5"},
		{"duplicate", "Key,Key", "Key,Key"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			info := BatchStatusInfo{Bucket: "bucket", ShowFields: tc.fields}
			require.Nil(t, info.Check())
			assert.Equal(t, tc.want, info.ShowFields)
		})
	}
	for _, fields := range []string{"Key,Unknown", "Key,,Expiration", "Key,"} {
		info := BatchStatusInfo{Bucket: "bucket", ShowFields: fields}
		require.NotNil(t, info.Check(), fields)
	}
}

// TestFormatBatchStatus 验证完整字段值、默认六列兼容、自选顺序和空值占位。
func TestFormatBatchStatus(t *testing.T) {
	status, restore := 1, 2
	expiration, ia, it, archive, deep := int64(1706000000), int64(1701000000), int64(1702000000), int64(1704000000), int64(1705000000)
	result := batch.OperationResult{
		FSize: 123, Hash: "etag", MimeType: "text/plain", PutTime: 17000000001234567, Type: 5,
		EndUser: "owner", MD5: "md5", Status: &status, RestoreStatus: &restore,
		Expiration: &expiration, TransitionToIA: &ia, TransitionToIntelligentTiering: &it,
		TransitionToArchiveIR: 1703000000, TransitionToArchive: &archive, TransitionToDeepArchive: &deep,
	}
	assert.Equal(t, "a.txt\t123\tetag\ttext/plain\t17000000001234567\t5", formatBatchStatus("a.txt", &result, batchStatusFields[:6]))
	assert.Equal(t, "a.txt\t123\tetag\ttext/plain\t17000000001234567\t5\towner\tmd5\t1\t2\t1706000000\t1701000000\t1702000000\t1703000000\t1704000000\t1705000000",
		formatBatchStatus("a.txt", &result, batchStatusFields))
	assert.Equal(t, "1706000000\ta.txt\towner", formatBatchStatus("a.txt", &result, []string{"Expiration", "Key", "EndUser"}))

	empty := batch.OperationResult{}
	assert.Equal(t, "\ta.txt\t\t0\t", formatBatchStatus("a.txt", &empty, []string{"EndUser", "Key", "Expiration", "TransitionToArchiveIR", "MD5"}))
	assert.Equal(t, []string{"", "a.txt", ""}, strings.Split(formatBatchStatus("a.txt", &empty, []string{"Status", "Key", "Expiration"}), "\t"))
	assert.Len(t, strings.Split(formatBatchStatus("a.txt", &empty, batchStatusFields), "\t"), 16)
	assert.Equal(t, strings.Repeat("\t", 6), formatBatchStatus("a.txt", &empty,
		[]string{"Status", "RestoreStatus", "Expiration", "TransitionToIA", "TransitionToIntelligentTiering", "TransitionToArchive", "TransitionToDeepArchive"}))
	zero := int64(0)
	empty.Expiration = &zero
	assert.Equal(t, "0", formatBatchStatus("a.txt", &empty, []string{"Expiration"}))
}

// TestStatusResultDecoding 验证扩展公共批量结果后，单文件查询仍解码到原有的值类型字段。
func TestStatusResultDecoding(t *testing.T) {
	var result object.StatusResult
	require.NoError(t, json.Unmarshal([]byte(`{"status":1,"restoreStatus":2,"md5":"digest","expiration":1706000000,"transitionToARCHIVE":1704000000}`), &result))
	assert.Equal(t, 1, result.Status)
	assert.Equal(t, 2, result.RestoreStatus)
	assert.Equal(t, "digest", result.MD5)
	assert.Equal(t, int64(1706000000), result.Expiration)
	assert.Equal(t, int64(1704000000), result.TransitionToARCHIVE)
}
