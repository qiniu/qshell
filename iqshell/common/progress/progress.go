package progress

type Progress interface {
	Start()
	SetFileSize(fileSize int64)
	SendSize(newSize int64)
	Progress(current int64)
	End()
}