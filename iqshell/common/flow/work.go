package flow

type Work interface {
	WorkId() string
}

type WorkCreator interface {
	Create(info string)(work Work, err error)
}