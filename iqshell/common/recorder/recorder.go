package recorder

type Recorder interface {

	// Get 获取记录
	Get(key string) (value string, err error)

	// Put 添加记录
	Put(key, value string) error

	// Delete 删除记录
	Delete(key string) error
}
