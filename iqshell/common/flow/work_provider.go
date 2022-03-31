package flow

const UnknownWorkCount = int64(-1)

type WorkProvider interface {
	WorkTotalCount() int64
	Provide() (hasMore bool, work Work, err error)
}
