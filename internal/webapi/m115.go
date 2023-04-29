package webapi

type M115Response struct {
	State      bool   `json:"state"`
	ErrorCode  int    `json:"errcode,omitempty"`
	ErrorCode2 int    `json:"errno,omitempty"`
	ErrorMsg   string `json:"msg"`
	Data       string `json:"data"`
}

func (r *M115Response) Err() error {
	if !r.State {
		return getError(findNonZero(r.ErrorCode, r.ErrorCode2))
	}
	return nil
}
