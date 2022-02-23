package data

const (
	TrueString  = "true"
	FalseString = "false"

	ResumeApiV1 = "v1"
	ResumeApiV2 = "v2"

	BLOCK_BITS = 22              // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS // BLOCK SIZE
)

func GetBoolString(value bool) string {
	if value {
		return TrueString
	} else {
		return FalseString
	}
}

func Bool(value string) bool {
	return value == TrueString
}

const (
	StatusOK         = iota // process success
	StatusError             // process error
	StatusHalt              // local error
	StatusUserCancel        // 用户取消
)

// fetch 接口返回的结构
type FetchItem struct {
	RemoteUrl string
	Bucket    string
	Key       string
}

const DefaultLineSeparate = "\t"

type Checker interface {
	Check() error
}
