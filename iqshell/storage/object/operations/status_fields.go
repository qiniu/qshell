package operations

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

// batchStatusFields 列出批量查询支持的字段，前六项为默认输出顺序。
var batchStatusFields = []string{
	"Key", "Fsize", "Hash", "MimeType", "PutTime", "Type", "EndUser", "MD5", "Status", "RestoreStatus",
	"Expiration", "TransitionToIA", "TransitionToIntelligentTiering", "TransitionToArchiveIR", "TransitionToArchive", "TransitionToDeepArchive",
}

// formatBatchStatus 按指定顺序拼接 Tab 分隔的查询结果，空指针和空字符串均保留空列。
func formatBatchStatus(key string, result *batch.OperationResult, fields []string) string {
	values := map[string]interface{}{
		"Key":                            key,
		"Fsize":                          result.FSize,
		"Hash":                           result.Hash,
		"MimeType":                       result.MimeType,
		"PutTime":                        result.PutTime,
		"Type":                           result.Type,
		"EndUser":                        result.EndUser,
		"MD5":                            result.MD5,
		"Status":                         result.Status,
		"RestoreStatus":                  result.RestoreStatus,
		"Expiration":                     result.Expiration,
		"TransitionToIA":                 result.TransitionToIA,
		"TransitionToIntelligentTiering": result.TransitionToIntelligentTiering,
		"TransitionToArchiveIR":          result.TransitionToArchiveIR,
		"TransitionToArchive":            result.TransitionToArchive,
		"TransitionToDeepArchive":        result.TransitionToDeepArchive,
	}
	items := make([]string, len(fields))
	for i, field := range fields {
		switch value := values[field].(type) {
		case *int:
			if value != nil {
				items[i] = strconv.Itoa(*value)
			}
		case *int64:
			if value != nil {
				items[i] = strconv.FormatInt(*value, 10)
			}
		default:
			items[i] = fmt.Sprint(value)
		}
	}
	return strings.Join(items, "\t")
}
