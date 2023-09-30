package api

import (
	"encoding/json"
	"strings"

	"github.com/deadblue/elevengo/internal/apibase"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

type ShareDuration int

const (
	ShareOneDay  ShareDuration = 1
	ShareOneWeek ShareDuration = 7
	ShareForever ShareDuration = -1

	ShareStateAuditing = 0
	ShareStateAccepted = 1
	ShareStateRejected = 6
)

type ShareInfo struct {
	ShareCode string `json:"share_code"`

	ShareState    util.IntNumber `json:"share_state"`
	ShareTitle    string         `json:"share_title"`
	ShareUrl      string         `json:"share_url"`
	ShareDuration util.IntNumber `json:"share_ex_time"`
	ReceiveCode   string         `json:"receive_code"`

	ReceiveCount util.IntNumber `json:"receive_count"`

	FileCount   int            `json:"file_count"`
	FolderCount int            `json:"folder_count"`
	TotalSize   util.IntNumber `json:"total_size"`
}

type ShareListResult struct {
	Count int
	Items []*ShareInfo
}

//lint:ignore U1000 This type is used in generic.
type _ShareListResp struct {
	apibase.BasicResp

	Count int             `json:"count"`
	List  json.RawMessage `json:"list"`
}

func (r *_ShareListResp) Extract(v any) (err error) {
	ptr, ok := v.(*ShareListResult)
	if !ok {
		return errors.ErrUnsupportedResult
	}
	if err = json.Unmarshal(r.List, &ptr.Items); err != nil {
		return
	}
	ptr.Count = r.Count
	return
}

type ShareListSpec struct {
	apibase.JsonApiSpec[ShareListResult, _ShareListResp]
}

func (s *ShareListSpec) Init(offset int, userId string) *ShareListSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/slist")
	s.QuerySet("user_id", userId)
	s.QuerySetInt("offset", offset)
	s.QuerySetInt("limit", FileListLimit)
	return s
}

type ShareSendSpec struct {
	apibase.JsonApiSpec[ShareInfo, apibase.StandardResp]
}

func (s *ShareSendSpec) Init(fileIds []string, userId string) *ShareSendSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/send")
	s.FormSet("user_id", userId)
	s.FormSet("file_ids", strings.Join(fileIds, ","))
	s.FormSet("ignore_warn", "1")
	return s
}

type ShareGetSpec struct {
	apibase.JsonApiSpec[ShareInfo, apibase.StandardResp]
}

func (s *ShareGetSpec) Init(shareCode string) *ShareGetSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/shareinfo")
	s.QuerySet("share_code", shareCode)
	return s
}

type ShareUpdateSpec struct {
	apibase.VoidApiSpec
}

func (s *ShareUpdateSpec) Init(
	shareCode string, receiveCode string, duration ShareDuration,
) *ShareUpdateSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/share/updateshare")
	s.FormSet("share_code", shareCode)
	if receiveCode == "" {
		s.FormSet("auto_fill_recvcode", "1")
	} else {
		s.FormSet("receive_code", receiveCode)
	}
	if duration > 0 {
		s.FormSetInt("share_duration", int(duration))
	}
	return s
}

type ShareCancelSpec struct {
	apibase.VoidApiSpec
}

func (s *ShareCancelSpec) Init(shareCode string) *ShareCancelSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/share/updateshare")
	s.FormSet("share_code", shareCode)
	s.FormSet("action", "cancel")
	return s
}
