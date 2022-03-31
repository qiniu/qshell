package flow

type WorkProvider interface {
	Provide() (hasMore bool, work Work, err error)
}