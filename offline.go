package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"time"
)

const (
	apiOfflineSpace = "https://115.com/"
)

type ClearParams struct {
	flag int
}

func (cp *ClearParams) All(delete bool) *ClearParams {
	if delete {
		cp.flag = 5
	} else {
		cp.flag = 1
	}
	return cp
}
func (cp *ClearParams) Complete(delete bool) *ClearParams {
	if delete {
		cp.flag = 4
	} else {
		cp.flag = 0
	}
	return cp
}
func (cp *ClearParams) Failed() *ClearParams {
	cp.flag = 2
	return cp
}
func (cp *ClearParams) Running() *ClearParams {
	cp.flag = 3
	return cp
}

func (c *Client) offlineSpace() (err error) {
	qs := core.NewQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithInt64("_", time.Now().Unix())
	result := &internal.OfflineSpaceResult{}
	if err = c.hc.JsonApi(apiOfflineSpace, qs, nil, result); err != nil {
		return
	}
	// store to client
	if c.ot == nil {
		c.ot = &internal.OfflineToken{}
	}
	c.ot.Sign = result.Sign
	c.ot.Time = result.Time
	return nil
}

func (c *Client) callOfflineApi(url string, form core.Form, result interface{}) (err error) {
	if c.ot == nil {
		if err = c.offlineSpace(); err != nil {
			return
		}
	}
	if form == nil {
		form = core.NewForm()
	}
	form.WithString("uid", c.ui.UserId).
		WithString("sign", c.ot.Sign).
		WithInt64("time", c.ot.Time)
	err = c.hc.JsonApi(url, nil, form, result)
	return
}

//
//func (c *Client) OfflineList(page int) (tasks []*OfflineTask, remain int, err error) {
//	form := core.NewForm().WithInt("page", page)
//	result := &_OfflineListResult{}
//	err = c.callOfflineApi(apiOfflineList, form, result)
//	if err == nil {
//		if !result.State {
//			err = apiError(result.ErrorNo)
//		} else {
//			c.offline.QuotaTotal = result.QuotaTotal
//			c.offline.QuotaRemain = result.Quota
//			tasks, remain = result.Tasks, result.PageCount-1
//		}
//	}
//	return
//}
//
//func (c *Client) OfflineDelete(hash ...string) (err error) {
//	form := core.NewForm().WithStrings("hash", hash)
//	result := &_OfflineBasicResult{}
//	err = c.callOfflineApi(apiOfflineDelete, form, result)
//	if err == nil && !result.State {
//		err = apiError(result.ErrorNo)
//	}
//	return
//}
//
//func (c *Client) OfflineClear(flag ClearFlag) (err error) {
//	form := core.NewForm().WithInt("flag", int(flag))
//	result := &_OfflineBasicResult{}
//	err = c.callOfflineApi(apiOfflineClear, form, result)
//	if err == nil && !result.State {
//		err = apiError(result.ErrorNo)
//	}
//	return
//}
//
//func (c *Client) OfflineAddUrl(url string) (hash string, err error) {
//	form := core.NewForm().
//		WithString("url", url)
//	result := &_OfflineAddResult{}
//	err = c.callOfflineApi(apiOfflineAddUrl, form, result)
//	if err == nil {
//		if !result.State {
//			err = apiError(result.ErrorNo)
//		} else {
//			hash = result.InfoHash
//		}
//	}
//	return
//}
//
//type TorrentFileFilter func(path string, size int64) bool
//
//func (c *Client) OfflineAddTorrent(torrentFile string, filter TorrentFileFilter) (hash string, err error) {
//	// get torrent dir
//	qs := core.NewQueryString().
//		WithString("ct", "lixian").
//		WithString("ac", "get_id").
//		WithString("torrent", "1").
//		WithInt64("_", time.Now().Unix())
//	gdr := &_OfflineGetDirResult{}
//	if err = c.requestJson(apiBasic, qs, nil, gdr); err != nil {
//		return
//	}
//	// upload torrent
//	cf, err := c.UploadFile(gdr.CategoryId, torrentFile, "")
//	if err != nil {
//		return
//	}
//	// get torrent info
//	form := core.NewForm().
//		WithString("pickcode", cf.PickCode).
//		WithString("sha1", cf.Sha1)
//	tir := &_OfflineTorrentInfoResult{}
//	if err = c.callOfflineApi(apiOfflineTorrentInfo, form, tir); err != nil {
//		return
//	}
//	// add bt task
//	wanted, selectCount := make([]string, tir.FileCount), 0
//	for index, tf := range tir.FileList {
//		if filter == nil || filter(tf.Path, tf.Size) {
//			wanted[selectCount] = strconv.Itoa(index)
//			selectCount += 1
//		}
//	}
//	if selectCount == 0 {
//		return "", ErrOfflineNothindToAdd
//	}
//	form = core.NewForm().
//		WithString("savepath", tir.TorrentName).
//		WithString("info_hash", tir.InfoHash).
//		WithString("wanted", strings.Join(wanted[:selectCount], ","))
//	result := &_OfflineAddResult{}
//	err = c.callOfflineApi(apiOfflineAddTorrent, form, result)
//	if err == nil {
//		if !result.State {
//			err = apiError(result.ErrorNo)
//		} else {
//			hash = result.InfoHash
//		}
//	}
//	return
//}
