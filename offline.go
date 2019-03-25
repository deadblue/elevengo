package elevengo

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

func (c *Client) OfflineList(page int) (result *OfflineListResult, err error) {
	result = &OfflineListResult{}
	form := newForm(false).WithInt("page", page)
	if err = c.callOfflineApi(apiOfflineList, form, result); err != nil {
		return nil, err
	}
	return
}

func (c *Client) OfflineAddUrls(url ...string) (tasks []*OfflineAddUrlResult, err error) {
	form := newForm(false)
	if len(url) == 1 {
		form.WithString("url", url[0])
		result := &OfflineAddUrlResult{}
		if err = c.callOfflineApi(apiOfflineAddUrl, form, result); err != nil {
			return
		}
		tasks = []*OfflineAddUrlResult{result}
	} else {
		form.WithStrings("url", url)
		result := &OfflineAddUrlsResult{}
		if err = c.callOfflineApi(apiOfflineAddUrls, form, result); err != nil {
			return
		}
		tasks = result.Result
	}
	return
}

func (c *Client) OfflineDelete(hash ...string) (err error) {
	form := newForm(false).WithStrings("hash", hash)
	result := &OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineDelete, form, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorCode)
	}
	return
}

func (c *Client) OfflineClear(flag ClearFlag) (err error) {
	form := newForm(false).WithInt("flag", int(flag))
	result := &OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineClear, form, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorCode)
	}
	return
}
