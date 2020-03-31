package internal

type CaptchaSignResult struct {
	State bool   `json:"state"`
	Sign  string `json:"sign"`
}

type CaptchaSubmitResult struct {
	State bool `json:"state"`
}
