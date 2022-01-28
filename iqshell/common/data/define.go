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

const (
	STATUS_OK    = iota // process success
	STATUS_ERROR        //process error
	STATUS_HALT         //local error
)

// fetch 接口返回的结构
type FetchItem struct {
	RemoteUrl string
	Bucket    string
	Key       string
}

const DefaultLineSeparate = "\t"
