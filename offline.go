package elevengo

import "fmt"

func (c *Client) updateOfflineSpace() (err error) {
	qs := newQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithTimestamp("_")
	result := &OfflineSpaceResult{}
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

func (c *Client) OfflineTaskList(page int) (result *OfflineTaskListResult, err error) {
	result = &OfflineTaskListResult{}
	form := newForm(false).WithInt("page", page)
	if err = c.callOfflineApi(apiOfflineList, form, result); err != nil {
		return nil, err
	}
	return
}

func (c *Client) OfflineTaskAddUrls(url ...string) (tasks []*OfflineTaskAddUrlResult, err error) {
	form := newForm(false)
	if len(url) == 1 {
		form.WithString("url", url[0])
		result := &OfflineTaskAddUrlResult{}
		if err = c.callOfflineApi(apiOfflineAddUrl, form, result); err != nil {
			return
		}
		tasks = []*OfflineTaskAddUrlResult{result}
	} else {
		form.WithStrings("url", url)
		result := &OfflineTaskAddUrlsResult{}
		if err = c.callOfflineApi(apiOfflineAddUrls, form, result); err != nil {
			return
		}
		tasks = result.Result
	}
	return
}

func (c *Client) OfflineTaskDelete(hash ...string) (err error) {
	form := newForm(false).WithStrings("hash", hash)
	result := &OfflineBasicResult{}
	if err = c.callOfflineApi(apiOfflineDelete, form, result); err != nil {
		return
	}
	if !result.State {
		err = fmt.Errorf("api error: %d", result.ErrorCode)
	}
	return
}

func (c *Client) OfflineTaskClear(flag ClearFlag) (err error) {
	form := newForm(false).WithInt("flag", int(flag))
	result := &OfflineBasicResult{}
	if err = c.callOfflineApi(apiOfflineClear, form, result); err != nil {
		return err
	}
	if !result.State {
		err = fmt.Errorf("api error: %d", result.ErrorCode)
	}
	return err
}
