package elevengo

// Get file list under the specific category.
//
// "0" is a special categoryId which means the root,
// everything starts from here.
//
// `offset` is base on zero.
//
// `limit` can not be lower than `FileListMinLimit`,
//  and can not be higher than `FileListMaxLimit`
//
// `sort` is optional, pass `nil` will use the default sort option:
// sorting by modify time in desc.
func (c *Client) FileList(categoryId string, offset, limit int, sort *SortOptions) (result *FileListResult, err error) {
	if limit < FileListMinLimit {
		limit = FileListMinLimit
	} else if limit > FileListMaxLimit {
		limit = FileListMaxLimit
	}
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("o", string(OrderByTime)).
		WithString("asc", "0").
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithInt("offset", offset).
		WithInt("limit", limit)
	// override default sort parameters
	if sort != nil {
		if sort.Asc {
			qs.WithString("asc", "1")
		}
		if sort.Flag != nil {
			qs.WithString("o", string(*sort.Flag))
		}
	}
	result = &FileListResult{}
	err = c.requestJson(apiFileList, qs, nil, result)
	if err == nil && !result.State {
		err = apiError(result.MessageCode)
	}
	return
}

// Search files which's name contains the specific keyword,
// the searching is recursive, starts from the specific category.
//
// `keyword` can not be empty
//
// `offset` is base on zero.
//
// `limit` can not be lower than `FileListMinLimit`,
//  and can not be higher than `FileListMaxLimit`
func (c *Client) FileSearch(categoryId, keyword string, offset, limit int) (result *FileSearchResult, err error) {
	if len(keyword) == 0 {
		return nil, ErrEmptyKeyword
	}
	if limit < FileListMinLimit {
		limit = FileListMinLimit
	} else if limit > FileListMaxLimit {
		limit = FileListMaxLimit
	}
	qs := newQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("search_value", keyword).
		WithString("format", "json").
		WithInt("offset", offset).
		WithInt("limit", limit)
	result = &FileSearchResult{}
	err = c.requestJson(apiFileSearch, qs, nil, result)
	if err == nil && !result.State {
		err = apiError(result.MessageCode)
	}
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
