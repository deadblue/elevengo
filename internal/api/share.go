package api

import (
	"encoding/json"
	"strings"

	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
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

	ShareState    base.IntNumber `json:"share_state"`
	ShareTitle    string         `json:"share_title"`
	ShareUrl      string         `json:"share_url"`
	ShareDuration base.IntNumber `json:"share_ex_time"`
	ReceiveCode   string         `json:"receive_code"`

	ReceiveCount base.IntNumber `json:"receive_count"`

	FileCount   int            `json:"file_count"`
	FolderCount int            `json:"folder_count"`
	TotalSize   base.IntNumber `json:"total_size"`
}

type ShareListResult struct {
	Count int
	Items []*ShareInfo
}

//lint:ignore U1000 This type is used in generic.
type _ShareListResp struct {
	base.BasicResp

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
	base.JsonApiSpec[ShareListResult, _ShareListResp]
}

func (s *ShareListSpec) Init(offset int, userId string) *ShareListSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/slist")
	s.QuerySet("user_id", userId)
	s.QuerySetInt("offset", offset)
	s.QuerySetInt("limit", FileListLimit)
	return s
}

type ShareSendSpec struct {
	base.JsonApiSpec[ShareInfo, base.StandardResp]
}

func (s *ShareSendSpec) Init(fileIds []string, userId string) *ShareSendSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/send")
	s.FormSet("user_id", userId)
	s.FormSet("file_ids", strings.Join(fileIds, ","))
	s.FormSet("ignore_warn", "1")
	return s
}

type ShareGetSpec struct {
	base.JsonApiSpec[ShareInfo, base.StandardResp]
}

func (s *ShareGetSpec) Init(shareCode string) *ShareGetSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/shareinfo")
	s.QuerySet("share_code", shareCode)
	return s
}

type ShareUpdateSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *ShareUpdateSpec) Init(
	shareCode string, receiveCode string, duration ShareDuration,
) *ShareUpdateSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/updateshare")
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
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *ShareCancelSpec) Init(shareCode string) *ShareCancelSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/share/updateshare")
	s.FormSet("share_code", shareCode)
	s.FormSet("action", "cancel")
	return s
}
