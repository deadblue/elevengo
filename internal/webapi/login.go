package webapi

import "encoding/json"

type LoginBasicResponse struct {
	State   int    `json:"state"`
	Code    int    `json:"code"`
	Message string `json:"message"`

	ErrorCode    int    `json:"errno"`
	ErrorMessage string `json:"error"`

	Data json.RawMessage `json:"data"`
}

func (r *LoginBasicResponse) Err() error {
	if r.State != 0 {
		return nil
	}
	return getError(r.Code)
}

func (r *LoginBasicResponse) Decode(data interface{}) error {
	return json.Unmarshal(r.Data, data)
}

type LoginKeyData struct {
	Key string `json:"key"`
}

type LoginResultData struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}
