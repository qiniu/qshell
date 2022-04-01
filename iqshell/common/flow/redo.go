package flow

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type Redo interface {

	// ShouldRedo
	// @Description: 是否需要重新做
	// @param work 工作信息
	// @param workRecord 此工作的记录
	// @return shouldRedo 是否需要重做
	// @return cause 需要重做或不能重做的原因
	ShouldRedo(work Work, workRecord *WorkRecord) (shouldRedo bool, cause *data.CodeError)
}
