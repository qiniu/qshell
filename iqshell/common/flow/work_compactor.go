package flow

type WorkCompactor interface {
	Compact(work Work)(info string, err error)
}

