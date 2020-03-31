package elevengo

import (
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"strconv"
	"time"
)

const (
	apiFileIndex = "https://webapi.115.com/files/index_info"

	apiFileList       = "https://webapi.115.com/files"
	apiFileListByName = "https://aps.115.com/natsort/files.php"

	apiFileAdd    = "https://webapi.115.com/files/add"
	apiFileCopy   = "https://webapi.115.com/files/copy"
	apiFileMove   = "https://webapi.115.com/files/move"
	apiFileRename = "https://webapi.115.com/files/batch_rename" // params: files_new_name
	apiFileDelete = "https://webapi.115.com/rb/delete"
)

func (c *Client) FileIndex() (err error) {
	result := new(internal.FileIndexResult)
	err = c.hc.JsonApi(apiFileIndex, nil, nil, result)
	if err != nil {
		// TODO: handle api result
	}
	return
}

func parseTimestamp(s string) time.Time {
	sec, _ := strconv.ParseInt(s, 10, 64)
	return time.Unix(sec, 0)
}

func (c *Client) FileList(categoryId string, page *PageParam, sort *SortParam) (files []FileItem, err error) {
	// prepare parameters
	if sort == nil {
		sort = (&SortParam{}).ByTime().Desc()
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithInt("offset", page.offset()).
		WithInt("limit", page.limit()).
		WithString("o", sort.flag)
	if sort.asc {
		qs.WithString("asc", "1")
	} else {
		qs.WithString("asc", "0")
	}
	// select API URL
	apiUrl := apiFileList
	if sort.flag == "file_name" {
		apiUrl = apiFileListByName
	}
	// call API
	result := new(internal.FileListResult)
	err = c.hc.JsonApi(apiUrl, qs, nil, result)
	if err == nil && !result.State {
		err = fmt.Errorf("get file list failed")
	}
	if err != nil {
		return
	}
	// convert API result
	files = make([]FileItem, len(result.Data))
	for i, data := range result.Data {
		fi := FileItem{
			CategoryId: data.CategoryId,
			Name:       data.Name,
			PickCode:   data.PickCode,
			CreateTime: parseTimestamp(data.CreateTime),
			UpdateTime: parseTimestamp(data.UpdateTime),
		}
		if data.FileId != "" {
			fi.IsCategory = false
			fi.FileId = data.FileId
			fi.Size = data.Size
			fi.Sha1 = data.Sha1
		} else {
			fi.IsCategory = true
			fi.ParentId = data.ParentId
		}

		files[i] = fi
	}
	return
}

func (c *Client) FileCopy(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileCopy, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c Client) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileMove, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c *Client) FileDelete(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileDelete, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c *Client) CategoryAdd(parentId, name string) (categoryId string, err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &internal.CategoryAddResult{}
	err = c.hc.JsonApi(apiFileAdd, nil, form, result)
	if err != nil {
		return
	}
	if !result.State {
		// TODO: convert error
		err = errors.New(result.Error)
	} else {
		categoryId = result.CategoryId
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
//func (c *Client) FileSearch(categoryId, keyword string, page *PageParam) (files []*CloudFile, remain int, err error) {
//	if len(keyword) == 0 {
//		return nil, 0, ErrEmptyKeyword
//	}
//	if page == nil {
//		page = &PageParam{}
//	}
//	qs := core.NewQueryString().
//		WithString("aid", "1").
//		WithString("cid", categoryId).
//		WithString("search_value", keyword).
//		WithString("format", "json").
//		WithInt("offset", page.offset()).
//		WithInt("limit", page.limit())
//	// call API
//	result := &_FileSearchResult{}
//	err = c.requestJson(apiFileSearch, qs, nil, result)
//	if err == nil && !result.State {
//		err = apiError(result.MessageCode)
//	}
//	if err != nil {
//		return
//	}
//	// remain file count
//	remain = result.TotalCount - (result.Offset + result.PageSize)
//	if remain < 0 {
//		remain = 0
//	}
//	// convert result
//	files = make([]*CloudFile, len(result.Data))
//	for index, data := range result.Data {
//		info := &CloudFile{
//			IsCategory: false,
//			IsSystem:   false,
//			CategoryId: data.CategoryId,
//			Name:       data.Name,
//			Size:       data.Size,
//			PickCode:   data.PickCode,
//		}
//		info.CreateTime, _ = strconv.ParseInt(data.CreateTime, 10, 64)
//		info.UpdateTime, _ = strconv.ParseInt(data.UpdateTime, 10, 64)
//		if data.FileId == nil {
//			info.IsCategory = true
//			info.ParentId = *data.ParentId
//		} else {
//			info.FileId = *data.FileId
//			info.Sha1 = *data.Sha1
//		}
//		files[index] = info
//	}
//	return
//}

//func (c *Client) FileRename(fileId, name string) (err error) {
//	form := core.NewForm().
//		WithString("fid", fileId).
//		WithString("file_name", name)
//	result := new(_FileOperateResult)
//	err = c.requestJson(apiFileEdit, nil, form, result)
//	if err == nil && !result.State {
//		err = apiError(-1)
//	}
//	return
//}
//
