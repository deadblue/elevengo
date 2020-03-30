package elevengo

const Version = "0.0.1"

const (
	FileListMinLimit = 10
	FileListMaxLimit = 1000

	ErrorAccountNeedVerify  = 911
	ErrorOfflineIllegalTask = 10003
	ErrorFileNotExists      = 70004
	ErrorFileIncomplete     = 70005
)

const (
	domain = ".115.com"

	apiBasic = "https://115.com/"

	apiFileList       = "https://webapi.115.com/files"
	apiFileAdd        = "https://webapi.115.com/files/add"
	apiFileMove       = "https://webapi.115.com/files/move"
	apiFileCopy       = "https://webapi.115.com/files/copy"
	apiFileEdit       = "https://webapi.115.com/files/edit"
	apiFileDelete     = "https://webapi.115.com/rb/delete"
	apiFileSearch     = "https://webapi.115.com/files/search"
	apiFileDownload   = "https://webapi.115.com/files/download"
	apiFileVideo      = "https://webapi.115.com/files/video"
	apiCategoryGet    = "https://webapi.115.com/category/get"
	apiFileUploadInit = "https://uplb.115.com/3.0/sampleinitupload.php"

	apiOfflineList        = "https://115.com/web/lixian/?ct=lixian&ac=task_lists"
	apiOfflineDelete      = "https://115.com/web/lixian/?ct=lixian&ac=task_del"
	apiOfflineClear       = "https://115.com/web/lixian/?ct=lixian&ac=task_clear"
	apiOfflineAddUrl      = "https://115.com/web/lixian/?ct=lixian&ac=add_task_url"
	apiOfflineTorrentInfo = "https://115.com/web/lixian/?ct=lixian&ac=torrent"
	apiOfflineAddTorrent  = "https://115.com/web/lixian/?ct=lixian&ac=add_task_bt"

	apiCaptcha       = "https://captchaapi.115.com/"
	apiCaptchaSubmit = "https://webapi.115.com/user/captcha"

	apiQrcodeToken  = "https://qrcodeapi.115.com/api/1.0/web/1.0/token"
	apiQrcodeStatus = "https://qrcodeapi.115.com/get/status/"
	apiQrcodeLogin  = "https://passportapi.115.com/app/1.0/web/1.0/login/qrcode"

	qrcodeStatusWaiting int = iota
	qrcodeStatusScanned
	qrcodeStatusConfirmed
)
