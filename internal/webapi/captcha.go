package webapi

type CaptchaSignResponse struct {
	BasicResponse
	Sign string `json:"sign"`
}
