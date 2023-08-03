// Deprecated: "webapi" package will be replaced by "api" package.
package webapi

const (
	ApiUserInfo  = "https://my.115.com/?ct=ajax&ac=nav"
	ApiIndexInfo = "https://webapi.115.com/files/index_info"

	ApiLoginCheck = "https://passportapi.115.com/app/1.0/web/1.0/check/sso"

	ApiQrcodeStatus = "https://qrcodeapi.115.com/get/status/"

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
	ApiFileVideo = "https://webapi.115.com/files/video"

	ApiTaskList    = "https://lixian.115.com/lixian/?ct=lixian&ac=task_lists"
	ApiTaskDelete  = "https://lixian.115.com/lixian/?ct=lixian&ac=task_del"
	ApiTaskClear   = "https://lixian.115.com/lixian/?ct=lixian&ac=task_clear"
	ApiTaskAddUrls = "https://lixian.115.com/lixianssp/?ac=add_task_urls"

	ApiDownloadGetUrl = "https://proapi.115.com/app/chrome/downurl"

	ApiUploadInfo     = "https://proapi.115.com/app/uploadinfo"
	ApiUploadOssToken = "https://uplb.115.com/3.0/gettoken.php"
	ApiUploadInit     = "https://uplb.115.com/4.0/initupload.php"

	ApiUploadSimpleInit = "https://uplb.115.com/3.0/sampleinitupload.php"

	ApiGetVersion = "https://appversion.115.com/1/web/1.0/api/chrome"
)
