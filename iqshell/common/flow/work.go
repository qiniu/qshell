package flow

type Work interface {
	WorkId() string
}

type WorkInfo struct {
	Data string `json:"data"`
	Work Work   `json:"work"`
}
