package elevengo

import (
	"strconv"
	"strings"
)

func (c *Client) updateOfflineSpace() (err error) {
	qs := newQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithTimestamp("_")
	result := &_OfflineSpaceResult{}
	if err = c.requestJson(apiBasic, qs, nil, result); err != nil {
		return
	}
	// store to client
	if c.offline == nil {
		c.offline = &_OfflineToken{}
	}
	c.offline.Sign = result.Sign
	c.offline.Time = result.Time
	return nil
}

func (c *Client) callOfflineApi(url string, form *_Form, result interface{}) (err error) {
	if c.offline == nil {
		if err = c.updateOfflineSpace(); err != nil {
			return
		}
	}
	if form == nil {
		form = newForm(false)
	}
	form.WithString("uid", c.info.UserId).
		WithString("sign", c.offline.Sign).
		WithInt64("time", c.offline.Time)
	err = c.requestJson(url, nil, form, result)
	return
}

func (c *Client) OfflineList(page int) (result *OfflineListResult, err error) {
	result = &OfflineListResult{}
	form := newForm(false).WithInt("page", page)
	if err = c.callOfflineApi(apiOfflineList, form, result); err != nil {
		return nil, err
	}
	return
}

func (c *Client) OfflineDelete(hash ...string) (err error) {
	form := newForm(false).WithStrings("hash", hash)
	result := &_OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineDelete, form, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorCode)
	}
	return
}

func (c *Client) OfflineClear(flag ClearFlag) (err error) {
	form := newForm(false).WithInt("flag", int(flag))
	result := &_OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineClear, form, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorCode)
	}
	return
}

func (c *Client) OfflineAddUrls(url ...string) (tasks []*OfflineAddResult, err error) {
	form := newForm(false)
	if len(url) == 1 {
		form.WithString("url", url[0])
		result := &OfflineAddResult{}
		if err = c.callOfflineApi(apiOfflineAddUrl, form, result); err != nil {
			return
		}
		tasks = []*OfflineAddResult{result}
	} else {
		form.WithStrings("url", url)
		result := &_OfflineBatchAddResult{}
		if err = c.callOfflineApi(apiOfflineAddUrls, form, result); err != nil {
			return
		}
		tasks = result.Result
	}
	return
}

type TorrentFileFilter func(path string, size int64) bool

func (c *Client) OfflineAddBT(torrentFile string, filter TorrentFileFilter) (result *OfflineAddResult, err error) {
	// get torrent dir
	qs := newQueryString().
		WithString("ct", "lixian").
		WithString("ac", "get_id").
		WithString("torrent", "1").
		WithTimestamp("_")
	gdr := &_OfflineGetDirResult{}
	if err = c.requestJson(apiBasic, qs, nil, gdr); err != nil {
		return
	}
	// upload torrent
	uf, err := c.UploadFile(gdr.CategoryId, torrentFile, "")
	if err != nil {
		return
	}
	// get torrent info
	form := newForm(false).
		WithString("pickcode", uf.PickCode).
		WithString("sha1", uf.Sha1)
	tr := &_OfflineTorrentResult{}
	if err = c.callOfflineApi(apiOfflineTorrent, form, tr); err != nil {
		return
	}
	// add bt task
	wanted, selectCount := make([]string, tr.FileCount), 0
	for index, tf := range tr.FileList {
		if filter == nil || filter(tf.Path, tf.Size) {
			wanted[selectCount] = strconv.Itoa(index)
			selectCount += 1
		}
	}
	if selectCount == 0 {
		return nil, ErrOfflineNothindToAdd
	}
	form = newForm(false).
		WithString("savepath", tr.TorrentName).
		WithString("info_hash", tr.InfoHash).
		WithString("wanted", strings.Join(wanted[:selectCount], ","))
	result = &OfflineAddResult{}
	err = c.callOfflineApi(apiOfflineAddBt, form, result)
	if err != nil {
		return
	}
	return
}
