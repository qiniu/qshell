package qshell

//dir to store some cached files for qshell, like ak, sk
var QShellRootPath string

const (
	BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS
)
