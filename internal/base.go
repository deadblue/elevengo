package internal

type ApiResult interface {
	IsFailed() bool
}

type BaseApiResult struct {
	State bool `json:"state"`
}

func (r *BaseApiResult) IsFailed() bool {
	return !r.State
}
