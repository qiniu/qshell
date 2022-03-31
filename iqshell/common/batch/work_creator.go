package batch

type WorkCreatorInfo struct {
	EnableStdIn  bool   // 是否允许 stdin, 当 InputFile 不存在时使用 stdin
	ItemSeparate string // 分隔符
	InputFile    string // batch 操作输入文件
	Force        bool   // 无需验证即可 batch 操作，类似于二维码验证
}
