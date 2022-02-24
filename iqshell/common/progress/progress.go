package progress

type Progress interface {
	Start()
	Progress(total, current int64)
	End()
}