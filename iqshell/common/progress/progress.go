package progress

import "io"

type Progress interface {
	io.Writer

	Start()
	SetFileSize(fileSize int64)
	SendSize(newSize int64)
	Progress(current int64)
	End()
}
