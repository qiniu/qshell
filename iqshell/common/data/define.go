package data

const (
	TrueString  = "true"
	FalseString = "false"

	ResumeApiV1 = "v1"
	ResumeApiV2 = "v2"

	DefaultLineSeparate = "\t"

	BLOCK_BITS = 22              // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS // BLOCK SIZE
)

// 此处需要明确出值
const (
	StatusOK         = 0 // process success
	StatusError      = 1 // process error
	StatusUserCancel = 2 // 用户取消
	StatusHalt       = 3 // local error
)
