package flow

type WorkerProvider interface {
	Provide() (worker Worker, err error)
}
