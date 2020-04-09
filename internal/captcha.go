package internal

type CaptchaSignResult struct {
	BaseApiResult
	Sign      string `json:"sign"`
	ErrorCode int    `json:"errno"`
	Error     string `json:"error"`
}

type CaptchaSubmitResult struct {
	BaseApiResult
}
