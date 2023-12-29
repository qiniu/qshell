package object

import (
	"fmt"

	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
)

// ChangeLifecycleApiInfo
// 相关链接：https://developer.qiniu.com/kodo/8062/modify-object-life-cycle
type ChangeLifecycleApiInfo struct {
	Bucket                 string `json:"bucket"`
	Key                    string `json:"key"`
	ToIAAfterDays          int    `json:"to_ia_after_days"`           // 转换到 低频存储类型，设置为 -1 表示取消
	ToArchiveIRAfterDays   int    `json:"to_archive_ir_after_days"`   // 转换到 归档直读存储类型， 设置为 -1 表示取消
	ToArchiveAfterDays     int    `json:"to_archive_after_days"`      // 转换到 归档存储类型， 设置为 -1 表示取消
	ToDeepArchiveAfterDays int    `json:"to_deep_archive_after_days"` // 转换到 深度归档存储类型， 设置为 -1 表示取消
	DeleteAfterDays        int    `json:"delete_after_days"`          // 过期删除，删除后不可恢复，设置为 -1 表示取消
}

func (l *ChangeLifecycleApiInfo) GetBucket() string {
	return l.Bucket
}

func (l *ChangeLifecycleApiInfo) ToOperation() (string, *data.CodeError) {
	if len(l.Bucket) == 0 || len(l.Key) == 0 {
		return "", alert.CannotEmptyError("change lifecycle operation bucket or key", "")
	}

	lifecycleSetting := ""
	if l.ToIAAfterDays != 0 {
		lifecycleSetting += fmt.Sprintf("/toIAAfterDays/%d", l.ToIAAfterDays)
	}
	if l.ToArchiveIRAfterDays != 0 {
		lifecycleSetting += fmt.Sprintf("/toArchiveIRAfterDays/%d", l.ToArchiveIRAfterDays)
	}
	if l.ToArchiveAfterDays != 0 {
		lifecycleSetting += fmt.Sprintf("/toArchiveAfterDays/%d", l.ToArchiveAfterDays)
	}
	if l.ToDeepArchiveAfterDays != 0 {
		lifecycleSetting += fmt.Sprintf("/toDeepArchiveAfterDays/%d", l.ToDeepArchiveAfterDays)
	}
	if l.DeleteAfterDays != 0 {
		lifecycleSetting += fmt.Sprintf("/deleteAfterDays/%d", l.DeleteAfterDays)
	}

	if len(lifecycleSetting) == 0 {
		return "", data.NewEmptyError().AppendDesc("invalid change lifecycle operation, must set at least one value of lifecycle")
	}
	return fmt.Sprintf("/lifecycle/%s%s", storage.EncodedEntry(l.Bucket, l.Key), lifecycleSetting), nil
}

func (l *ChangeLifecycleApiInfo) WorkId() string {
	return fmt.Sprintf("ChangeLifecycle|%s|%s", l.Bucket, l.Key)
}

func ChangeLifecycle(info *ChangeLifecycleApiInfo) (*batch.OperationResult, *data.CodeError) {
	return batch.One(info)
}
