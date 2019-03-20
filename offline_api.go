package elevengo

func (c *Client) updateOfflineSpace() (err error) {
	qs := newQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithTimestamp("_")
	result := &OfflineSpaceResult{}
	if err = c.requestJson(apiHost, qs, nil, result); err != nil {
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

func (c *Client) callOfflineApi(action string, form *_Form, result interface{}) (err error) {
	if c.offline == nil {
		if err = c.updateOfflineSpace(); err != nil {
			return
		}
	}
	qs := newQueryString().
		WithString("ct", "lixian").
		WithString("ac", action)
	if form == nil {
		form = newForm(false)
	}
	form.WithString("uid", c.info.UserId).
		WithString("sign", c.offline.Sign).
		WithInt64("time", c.offline.Time)
	err = c.requestJson(apiOffline, qs, form, result)
	return
}

func (c *Client) OfflineTaskList(page int) (result *OfflineTaskListResult, err error) {
	result = &OfflineTaskListResult{}
	form := newForm(false).WithInt("page", page)
	if err = c.callOfflineApi(offlineActionTaskList, form, result); err != nil {
		return nil, err
	}
	return
}

func (c *Client) OfflineTaskDelete(hash ...string) (err error) {
	form := newForm(false).WithStrings("hash", hash)
	result := &OfflineBasicResult{}
	return c.callOfflineApi(offlineActionTaskDelete, form, result)
}

func (c *Client) OfflineTaskClear(flag ClearFlag) (err error) {
	form := newForm(false).WithInt("flag", int(flag))
	result := &OfflineBasicResult{}
	return c.callOfflineApi(offlineActionTaskClear, form, result)
}

func (c *Client) OfflineTaskAddUrls(url ...string) (err error) {
	action, form := "", newForm(false)
	var result interface{}
	if len(url) == 1 {
		action = offlineActionTaskAddUrl
		form.WithString("url", url[0])
		result = &OfflineTaskAddUrlResult{}
	} else {
		action = offlineActionTaskAddUrls
		form.WithStrings("url", url)
		result = &OfflineTaskAddUrlsResult{}
	}
	return c.callOfflineApi(action, form, result)
}
