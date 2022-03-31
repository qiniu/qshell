package flow

type WorkCreator interface {
	Create(info string) (work Work, err error)
}
