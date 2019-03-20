package elevengo

func (c *Client) FileList(dirId string, sort OrderMode) (result *FileListResult, err error) {
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", dirId).
		WithString("o", string(sort)).
		WithString("asc", "1").
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithInt("offset", 0).
		WithInt("limit", 100)
	result = &FileListResult{}
	err = c.requestJson(apiFileList, qs, nil, result)
	return
}

func (c *Client) FileMkdir(parentId, name string) (err error) {
	form := newForm(false).
		WithString("pid", parentId).
		WithString("cname", name)
	err = c.requestJson(apiFileAdd, nil, form, nil)
	return
}

func (c *Client) FileMove(parentId string, fileIds ...string) (err error) {
	form := newForm(false).
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	return c.requestJson(apiFileMove, nil, form, nil)
}

func (c *Client) FileRename(fileId, name string) (err error) {
	form := newForm(false).
		WithString("fid", fileId).
		WithString("file_name", name)
	return c.requestJson(apiFileEdit, nil, form, nil)
}

func (c *Client) FileDelete(parentId string, fileIds ...string) (err error) {
	form := newForm(false).
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	return c.requestJson(apiFileDelete, nil, form, nil)
}

func (c *Client) FileSearch(dirId, keyword string) (err error) {
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", dirId).
		WithString("search_value", keyword).
		WithString("format", "json").
		WithInt("offset", 0).
		WithInt("limit", 100)
	err = c.requestJson(apiFileSearch, qs, nil, nil)
	return
}

func (c *Client) FileDownload(pickcode string) (result *FileDownloadResult, err error) {
	qs := newQueryString().
		WithString("pickcode", pickcode).
		WithTimestamp("_")
	result = &FileDownloadResult{}
	err = c.requestJson(apiFileDownload, qs, nil, result)
	return
}
