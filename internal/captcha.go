package internal

type CaptchaSignResult struct {
	BaseApiResult
	Sign string `json:"sign"`
}

type CaptchaSubmitResult struct {
	BaseApiResult
}
