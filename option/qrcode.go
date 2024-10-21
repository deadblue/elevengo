package option

type QrcodeOptions struct {
	// App Type to login
	App string
}

func (o *QrcodeOptions) LoginWeb() *QrcodeOptions {
	o.App = "web"
	return o
}

func (o *QrcodeOptions) LoginAndroid() *QrcodeOptions {
	o.App = "android"
	return o
}

func (o *QrcodeOptions) LoginIos() *QrcodeOptions {
	o.App = "ios"
	return o
}

func (o *QrcodeOptions) LoginTv() *QrcodeOptions {
	o.App = "tv"
	return o
}

func (o *QrcodeOptions) LoginWechatMiniApp() *QrcodeOptions {
	o.App = "wechatmini"
	return o
}

func (o *QrcodeOptions) LoginAlipayMiniApp() *QrcodeOptions {
	o.App = "alipaymini"
	return o
}

func (o *QrcodeOptions) LoginQandroid() *QrcodeOptions {
	o.App = "qandroid"
	return o
}

func Qrcode() *QrcodeOptions {
	return (&QrcodeOptions{}).LoginWeb()
}
