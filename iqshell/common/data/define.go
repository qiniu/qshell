package data

const (
	TrueString  = "true"
	FalseString = "false"

	ResumeApiV1 = "v1"
	ResumeApiV2 = "v2"

	// Indicate that the blocksize is 4M
	BLOCK_BITS = 22
	// BLOCK SIZE
	BLOCK_SIZE = 1 << BLOCK_BITS
)

const (
	// process success
	STATUS_OK = iota
	//process error
	STATUS_ERROR
	//local error
	STATUS_HALT
)

// fetch 接口返回的结构
type FetchItem struct {
	RemoteUrl string
	Bucket    string
	Key       string
}
