package config

type Retry struct {
	Max      int `json:"max,omitempty"`
	Interval int `json:"interval,omitempty"`
}
