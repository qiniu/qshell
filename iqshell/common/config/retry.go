package config

type Retry struct {
	Max      int `json:"max,omitempty"`
	Interval int `json:"interval,omitempty"`
}

func (r *Retry) merge(from *Retry) {
	if from == nil {
		return
	}

	if r.Max == 0 {
		r.Max = from.Max
	}

	if r.Interval == 0 {
		r.Interval = from.Interval
	}
}
