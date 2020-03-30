package elevengo

//import (
//	"github.com/deadblue/elevengo/core"
//	"net/url"
//	"time"
//)
//
//func (c *Client) GetDownloadInfo(pickcode string) (info *DownloadInfo, err error) {
//	// call API to get download url
//	qs := core.NewQueryString().
//		WithString("pickcode", pickcode).
//		WithInt64("_", time.Now().Unix())
//	result := &_FileDownloadResult{}
//	err = c.requestJson(apiFileDownload, qs, nil, result)
//	if err == nil && !result.State {
//		err = apiError(result.ErrorNo)
//	}
//	if err != nil {
//		return
//	}
//	// get cookies for downloading
//	u, _ := url.Parse(result.FileUrl)
//	cookies := c.jar.Cookies(u)
//	// fill info
//	info = &DownloadInfo{
//		Url:       result.FileUrl,
//		UserAgent: c.ua,
//		Cookies:   make([]*DownloadCookie, len(cookies)),
//	}
//	for index, ck := range cookies {
//		info.Cookies[index] = &DownloadCookie{ck.Name, ck.Value}
//	}
//	return
//}
