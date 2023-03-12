package webapi

const (
	ApiUserInfo  = "https://my.115.com/?ct=ajax&ac=nav"
	ApiIndexInfo = "https://webapi.115.com/files/index_info"

	ApiLoginCheck = "https://passportapi.115.com/app/1.0/web/1.0/check/sso"

	ApiLoginGetKey   = "https://passportapi.115.com/app/1.0/web/5.0.1/login/getKey"
	ApiPasswordLogin = "https://passportapi.115.com/app/1.0/web/1.0/login/login"

	ApiSmsSendCode = "https://passportapi.115.com/app/1.0/web/1.0/code/sms/login"
	ApiSmsLogin    = "https://passportapi.115.com/app/1.0/web/1.0/login/vip"

	ApiQrcodeToken  = "https://qrcodeapi.115.com/api/1.0/web/1.0/token"
	ApiQrcodeStatus = "https://qrcodeapi.115.com/get/status/"
	ApiQrcodeLogin  = "https://passportapi.115.com/app/1.0/web/1.0/login/qrcode"

	ApiFileList       = "https://webapi.115.com/files"
	ApiFileListByName = "https://aps.115.com/natsort/files.php"
	ApiFileSearch     = "https://webapi.115.com/files/search"

	ApiFileInfo = "https://webapi.115.com/files/get_info"
	ApiFileStat = "https://webapi.115.com/category/get"

	ApiFileStar = "https://webapi.115.com/files/star"
	ApiFileEdit = "https://webapi.115.com/files/edit"

	ApiFileMove   = "https://webapi.115.com/files/move"
	ApiFileCopy   = "https://webapi.115.com/files/copy"
	ApiFileRename = "https://webapi.115.com/files/batch_rename"
	ApiFileDelete = "https://webapi.115.com/rb/delete"

	ApiFileFindDuplicate = "https://webapi.115.com/files/get_repeat_sha"

	ApiLabelList   = "https://webapi.115.com/label/list"
	ApiLabelAdd    = "https://webapi.115.com/label/add_multi"
	ApiLabelEdit   = "https://webapi.115.com/label/edit"
	ApiLabelDelete = "https://webapi.115.com/label/delete"

	ApiDirAdd      = "https://webapi.115.com/files/add"
	ApiDirGetId    = "https://webapi.115.com/files/getid"
	ApiDirSetOrder = "https://webapi.115.com/files/order"

	ApiFileImage = "https://webapi.115.com/files/image"
	ApiFileVideo = "https://v.anxia.com/webapi/files/video"

	ApiOfflineSpace   = "https://115.com/?ct=offline&ac=space"
	ApiOfflineList    = "https://115.com/web/lixian/?ct=lixian&ac=task_lists"
	ApiOfflineAddUrl  = "https://115.com/web/lixian/?ct=lixian&ac=add_task_url"
	ApiOfflineAddUrls = "https://115.com/web/lixian/?ct=lixian&ac=add_task_urls"
	ApiOfflineDelete  = "https://115.com/web/lixian/?ct=lixian&ac=task_del"
	ApiOfflineClear   = "https://115.com/web/lixian/?ct=lixian&ac=task_clear"

	ApiDownloadGetUrl = "https://proapi.115.com/app/chrome/downurl"

	ApiUploadInfo     = "https://proapi.115.com/app/uploadinfo"
	ApiUploadOssToken = "https://uplb.115.com/3.0/gettoken.php"
	ApiUploadInit     = "https://uplb.115.com/4.0/initupload.php"

	ApiUploadSimpleInit = "https://uplb.115.com/3.0/sampleinitupload.php"

	ApiCaptchaPage        = "https://captchaapi.115.com/?ac=security_code&type=web"
	ApiCaptchaCodeImage   = "https://captchaapi.115.com/?ct=index&ac=code&ctype=0"
	ApiCaptchaAllKeyImage = "https://captchaapi.115.com/?ct=index&ac=code&t=all"
	ApiCaptchaOneKeyImage = "https://captchaapi.115.com/?ct=index&ac=code&t=single"
	ApiCaptchaSign        = "https://captchaapi.115.com/?ac=code&t=sign"
	ApiCaptchaSubmit      = "https://webapi.115.com/user/captcha"

	ApiGetVersion = "https://appversion.115.com/1/web/1.0/api/chrome"
)
