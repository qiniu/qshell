package worker

type Flow interface {
	Run()
}

type FlowInfo struct {
	WorkerCount    int
	WorkProvider   WorkProvider
	WorkerProvider WorkerProvider
	GroupHandler   FlowHandler
}

func NewFlow(info FlowInfo) Flow {
	return &flow{
		info: info,
	}
}

type flow struct {
	info FlowInfo
}

func (g *flow) Run() {

}
