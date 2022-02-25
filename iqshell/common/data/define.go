package data

const (
	TrueString  = "true"
	FalseString = "false"

	ResumeApiV1 = "v1"
	ResumeApiV2 = "v2"

	StatusOK         = iota // process success
	StatusError             // process error
	StatusHalt              // local error
	StatusUserCancel        // 用户取消

	DefaultLineSeparate = "\t"

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

func BoolValue(value string) bool {
	return value == TrueString
}
