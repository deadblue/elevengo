package option

type QrcodeOption interface {
	isQrcodeOption()
}

type QrcodeLoginOption string

func (o QrcodeLoginOption) isQrcodeOption() {}

const (
	QrcodeLoginWeb     = QrcodeLoginOption("web")
	QrcodeLoginMac     = QrcodeLoginOption("mac")
	QrcodeLoginLinux   = QrcodeLoginOption("linux")
	QrcodeLoginWindows = QrcodeLoginOption("windows")
)
