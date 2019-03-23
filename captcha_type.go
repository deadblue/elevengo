package elevengo

type CaptchaSession struct {
	Callback  string
	CodeValue []byte
	CodeKeys  []byte
}

type _CaptchaSignResult struct {
	State bool   `json:"state"`
	Sign  string `json:"sign"`
}
