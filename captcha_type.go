package elevengo

type CaptchaSession struct {
	Callback  string
	CodeValue []byte
	CodeKeys  []byte
}
