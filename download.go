package elevengo

import "net/url"

func (c *Client) GetDownloadInfo(pickcode string) (info *DownloadInfo, err error) {
	// call API to get download url
	qs := newQueryString().
		WithString("pickcode", pickcode).
		WithTimestamp("_")
	result := &FileDownloadResult{}
	err = c.requestJson(apiFileDownload, qs, nil, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorNo)
	}
	if err != nil {
		return
	}
	// get cookies for downloading
	u, _ := url.Parse(result.FileUrl)
	cookies := c.jar.Cookies(u)
	// fill info
	info = &DownloadInfo{
		Url:       result.FileUrl,
		UserAgent: c.ua,
		Cookies:   make([]*DownloadCookie, len(cookies)),
	}
	for index, ck := range cookies {
		info.Cookies[index] = &DownloadCookie{ck.Name, ck.Value}
	}
	return
}
