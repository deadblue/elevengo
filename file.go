package elevengo

func (c *Client) FileList(categoryId string, sort OrderMode) (result *FileListResult, err error) {
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
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

func (c *Client) FileInfo(fileId string) (data *FileInfoData, err error) {
	form := newForm(false).WithString("file_id", fileId)
	result := &FileInfoResult{}
	if err = c.requestJson(apiFileInfo, nil, form, result); err == nil {
		if !result.State {
			err = apiError(result.ErrorNo)
		} else if len(result.Data) > 0 {
			data = result.Data[0]
		}
	}
	return
}

func (c *Client) FileCopy(parentId string, fileIds ...string) (err error) {
	form := newForm(false).
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	return c.requestJson(apiFileCopy, nil, form, nil)
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

func (c *Client) FileSearch(categoryId, keyword string) (err error) {
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("search_value", keyword).
		WithString("format", "json").
		WithInt("offset", 0).
		WithInt("limit", 100)
	err = c.requestJson(apiFileSearch, qs, nil, nil)
	return
}

func (c *Client) FileAddCategory(parentId, name string) (err error) {
	form := newForm(false).
		WithString("pid", parentId).
		WithString("cname", name)
	result := &FileAddResult{}
	if err = c.requestJson(apiFileAdd, nil, form, result); err == nil {
		if !result.State {
			err = apiError(result.ErrorNo)
		}
	}
	return
}

func (c *Client) FileGetCategory(categoryId string) (err error) {
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId)
	result := &CategoryGetResult{}
	return c.requestJson(apiCategoryGet, qs, nil, result)
}
