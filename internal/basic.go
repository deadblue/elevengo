package internal

type ApiResult interface {
	IsOK() bool
}

type BaseApiResult struct {
	State bool `json:"state"`
}

func (r *BaseApiResult) IsOK() bool {
	return r.State
}
