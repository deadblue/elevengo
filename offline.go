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

func (c *Client) OfflineList(page int) (tasks []*OfflineTask, remain int, err error) {
	form := newForm(false).WithInt("page", page)
	result := &_OfflineListResult{}
	err = c.callOfflineApi(apiOfflineList, form, result)
	if err == nil {
		if !result.State {
			err = apiError(result.ErrorNo)
		} else {
			c.offline.QuotaTotal = result.QuotaTotal
			c.offline.QuotaRemain = result.Quota
			tasks, remain = result.Tasks, result.PageCount-1
		}
	}
	return
}

func (c *Client) OfflineDelete(hash ...string) (err error) {
	form := newForm(false).WithStrings("hash", hash)
	result := &_OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineDelete, form, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorNo)
	}
	return
}

func (c *Client) OfflineClear(flag ClearFlag) (err error) {
	form := newForm(false).WithInt("flag", int(flag))
	result := &_OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineClear, form, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorNo)
	}
	return
}

func (c *Client) OfflineAddUrl(url string) (hash string, err error) {
	form := newForm(false).
		WithString("url", url)
	result := &_OfflineAddResult{}
	err = c.callOfflineApi(apiOfflineAddUrl, form, result)
	if err == nil {
		if !result.State {
			err = apiError(result.ErrorNo)
		} else {
			hash = result.InfoHash
		}
	}
	return
}

type TorrentFileFilter func(path string, size int64) bool

func (c *Client) OfflineAddTorrent(torrentFile string, filter TorrentFileFilter) (hash string, err error) {
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
	cf, err := c.UploadFile(gdr.CategoryId, torrentFile, "")
	if err != nil {
		return
	}
	// get torrent info
	form := newForm(false).
		WithString("pickcode", cf.PickCode).
		WithString("sha1", cf.Sha1)
	tir := &_OfflineTorrentInfoResult{}
	if err = c.callOfflineApi(apiOfflineTorrentInfo, form, tir); err != nil {
		return
	}
	// add bt task
	wanted, selectCount := make([]string, tir.FileCount), 0
	for index, tf := range tir.FileList {
		if filter == nil || filter(tf.Path, tf.Size) {
			wanted[selectCount] = strconv.Itoa(index)
			selectCount += 1
		}
	}
	if selectCount == 0 {
		return "", ErrOfflineNothindToAdd
	}
	form = newForm(false).
		WithString("savepath", tir.TorrentName).
		WithString("info_hash", tir.InfoHash).
		WithString("wanted", strings.Join(wanted[:selectCount], ","))
	result := &_OfflineAddResult{}
	err = c.callOfflineApi(apiOfflineAddTorrent, form, result)
	if err == nil {
		if !result.State {
			err = apiError(result.ErrorNo)
		} else {
			hash = result.InfoHash
		}
	}
	return
}
