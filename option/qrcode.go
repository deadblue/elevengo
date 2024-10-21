package option

type QrcodeOption interface {
	isQrcodeOption()
}

type QrcodeAppOption string

func (o QrcodeAppOption) isQrcodeOption() {}

const (
	QrcodeAppWeb      QrcodeAppOption = "web"
	QrcodeAppAndroid  QrcodeAppOption = "android"
	QrcodeAppIos      QrcodeAppOption = "ios"
	QrcodeAppTv       QrcodeAppOption = "tv"
	QrcodeAppAlipay   QrcodeAppOption = "alipaymini"
	QrcodeAppWechat   QrcodeAppOption = "wechatmini"
	QrcodeAppQandroid QrcodeAppOption = "qandroid"
)
