package elevengo

//
//type _UserInfo struct {
//	UserId string
//}
//
//type _OfflineToken struct {
//	Sign        string
//	Time        int64
//	QuotaTotal  int
//	QuotaRemain int
//}
//
//type _BasicResult struct {
//	State       bool    `json:"state"`
//	Error       *string `json:"error"`
//	ErrorType   *string `json:"errtype"`
//	MessageCode int     `json:"msg_code"`
//	Message     *string `json:"msg"`
//}
//
//type _FileRetrieveResult struct {
//	_BasicResult
//	ErrorNo int `json:"errNo"`
//}
//
//type _FileOperateResult struct {
//	_BasicResult
//	ErrorNo string `json:"errno"`
//}
//
//type _InnerFileListData struct {
//	ParentId   *string `json:"pid"`
//	CategoryId string  `json:"cid"`
//	FileId     *string `json:"fid"`
//	Name       string  `json:"n"`
//	Size       int64   `json:"s"`
//	PickCode   string  `json:"pc"`
//	Sha1       *string `json:"sha"`
//	CreateTime string  `json:"tp"`
//	UpdateTime string  `json:"te"`
//}
//
//type _FileListResult struct {
//	_FileRetrieveResult
//	TotalCount int                   `json:"count"`
//	SysCount   int                   `json:"sys_count"`
//	Offset     int                   `json:"offset"`
//	Limit      int                   `json:"limit"`
//	PageSize   int                   `json:"page_size"`
//	Data       []*_InnerFileListData `json:"data"`
//}
//
//type _FileSearchResult struct {
//	_FileRetrieveResult
//	TotalCount int                   `json:"count"`
//	Offset     int                   `json:"offset"`
//	PageSize   int                   `json:"page_size"`
//	Data       []*_InnerFileListData `json:"data"`
//}
//
//type _FileAddResult struct {
//	_FileOperateResult
//	AreaId       NumberString `json:"aid"`
//	CategoryId   NumberString `json:"cid"`
//	CategoryName string       `json:"cname"`
//	FileId       string       `json:"file_id"`
//	FileName     string       `json:"file_name"`
//}
//
//type _FileDownloadResult struct {
//	_FileRetrieveResult
//	FileId   string `json:"file_id"`
//	FileName string `json:"file_name"`
//	FileSize string `json:"file_size"`
//	Pickcode string `json:"pickcode"`
//	FileUrl  string `json:"file_url"`
//}
//
//type _FileUploadInitResult struct {
//	AccessKeyId string `json:"accessid"`
//	Callback    string `json:"callback"`
//	Expire      int    `json:"expire"`
//	UploadUrl   string `json:"host"`
//	ObjectKey   string `json:"object"`
//	Policy      string `json:"policy"`
//	Signature   string `json:"signature"`
//}
//
//type _InnerFileUploadData struct {
//	CategoryId string `json:"cid"`
//	FileId     string `json:"file_id"`
//	FileName   string `json:"file_name"`
//	FizeSize   string `json:"file_size"`
//	PickCode   string `json:"pick_code"`
//	Sha1       string `json:"sha1"`
//}
//
//type _FileUploadResult struct {
//	State   bool                  `json:"state"`
//	Code    int                   `json:"code"`
//	Message string                `json:"message"`
//	Data    *_InnerFileUploadData `json:"data"`
//}
//
//type _OfflineSpaceResult struct {
//	State bool    `json:"state"`
//	Data  float64 `json:"data"`
//	Size  string  `json:"size"`
//	Url   string  `json:"url"`
//	BtUrl string  `json:"bt_url"`
//	Limit int64   `json:"limit"`
//	Sign  string  `json:"sign"`
//	Time  int64   `json:"time"`
//}
//
//type _OfflineBasicResult struct {
//	State        bool    `json:"state"`
//	ErrorNo      int     `json:"errno"`
//	ErrorCode    int     `json:"errcode"`
//	ErrorType    string  `json:"errtype"`
//	ErrorMessage *string `json:"error_msg"`
//}
//
//type _OfflineListResult struct {
//	_OfflineBasicResult
//	Page       int            `json:"page"`
//	PageCount  int            `json:"page_count"`
//	PageRow    int            `json:"page_row"`
//	Count      int            `json:"count"`
//	Quota      int            `json:"quota"`
//	QuotaTotal int            `json:"total"`
//	Tasks      []*OfflineTask `json:"tasks"`
//}
//
//type _OfflineAddResult struct {
//	_OfflineBasicResult
//	InfoHash string `json:"info_hash"`
//	Name     string `json:"name"`
//}
//
//type _OfflineGetDirResult struct {
//	CategoryId string `json:"cid"`
//}
//
//type _TorrentFile struct {
//	Path string `json:"path"`
//	Size int64  `json:"size"`
//}
//
//type _OfflineTorrentInfoResult struct {
//	_OfflineBasicResult
//	TorrentName string          `json:"torrent_name"`
//	InfoHash    string          `json:"info_hash"`
//	FileSize    int64           `json:"file_size"`
//	FileCount   int             `json:"file_count"`
//	FileList    []*_TorrentFile `json:"torrent_filelist_web"`
//}
//
//type _CaptchaSignResult struct {
//	State bool   `json:"state"`
//	Sign  string `json:"sign"`
//}
