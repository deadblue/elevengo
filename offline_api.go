package elevengo

import "time"

func (c *Client) updateOfflineSpace() (err error) {
	qs := newRequestParameters().
		With("ct", "offline").
		With("ac", "space").
		WithInt64("_", time.Now().UnixNano())
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

func (c *Client) callOfflineApi(action string, body *RequestParameters, result interface{}) (err error) {
	if c.offline == nil {
		if err = c.updateOfflineSpace(); err != nil {
			return
		}
	}
	qs := newRequestParameters().
		With("ct", "lixian").
		With("ac", action)
	if body == nil {
		body = newRequestParameters()
	}
	body.With("uid", c.info.UserId).
		With("sign", c.offline.Sign).
		WithInt64("time", c.offline.Time)
	err = c.requestJson(apiOffline, qs, body.FormData(), result)
	return
}

func (c *Client) OfflineTaskList(page int) (result *OfflineTaskListResult, err error) {
	result = &OfflineTaskListResult{}
	body := newRequestParameters().WithInt("page", page)
	if err = c.callOfflineApi(offlineActionTaskList, body, result); err != nil {
		return nil, err
	}
	return
}

func (c *Client) OfflineTaskDelete(hash ...string) (err error) {
	body, result := newRequestParameters().WithStrings("hash", hash...), &OfflineBasicResult{}
	return c.callOfflineApi(offlineActionTaskDelete, body, result)
}

func (c *Client) OfflineTaskClear(flag ClearFlag) (err error) {
	body, result := newRequestParameters().WithInt("flag", int(flag)), &OfflineBasicResult{}
	return c.callOfflineApi(offlineActionTaskClear, body, result)
}

func (c *Client) OfflineTaskAddUrls(url ...string) (err error) {
	action, body := "", newRequestParameters()
	var result interface{}
	if len(url) == 1 {
		action = offlineActionTaskAddUrl
		body.With("url", url[0])
		result = &OfflineTaskAddUrlResult{}
	} else {
		action = offlineActionTaskAddUrls
		body.WithStrings("url", url...)
		result = &OfflineTaskAddUrlsResult{}
	}
	return c.callOfflineApi(action, body, result)
}
