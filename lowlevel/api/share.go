package api

import (
	"strings"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

type ShareDuration int

const (
	ShareOneDay  ShareDuration = 1
	ShareOneWeek ShareDuration = 7
	ShareForever ShareDuration = -1
)

type ShareListSpec struct {
	_JsonApiSpec[types.ShareListResult, protocol.ShareListResp]
}

func (s *ShareListSpec) Init(userId string, offset, limit int) *ShareListSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/share/slist")
	s.query.Set("user_id", userId).
		SetInt("offset", offset).
		SetInt("limit", limit)
	return s
}

type ShareSendSpec struct {
	_StandardApiSpec[types.ShareInfo]
}

func (s *ShareSendSpec) Init(fileIds []string, userId string) *ShareSendSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/share/send")
	s.form.Set("user_id", userId).
		Set("ignore_warn", "1").
		Set("file_ids", strings.Join(fileIds, ","))
	return s
}

type ShareGetSpec struct {
	_StandardApiSpec[types.ShareInfo]
}

func (s *ShareGetSpec) Init(shareCode string) *ShareGetSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/share/shareinfo")
	s.query.Set("share_code", shareCode)
	return s
}

type ShareUpdateSpec struct {
	_VoidApiSpec
}

func (s *ShareUpdateSpec) Init(
	shareCode string, receiveCode string, duration ShareDuration,
) *ShareUpdateSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/share/updateshare")
	s.form.Set("share_code", shareCode)
	if receiveCode == "" {
		s.form.Set("auto_fill_recvcode", "1")
	} else {
		s.form.Set("receive_code", receiveCode)
	}
	if duration > 0 {
		s.form.SetInt("share_duration", int(duration))
	}
	return s
}

type ShareCancelSpec struct {
	_VoidApiSpec
}

func (s *ShareCancelSpec) Init(shareCode string) *ShareCancelSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/share/updateshare")
	s.form.Set("share_code", shareCode)
	s.form.Set("action", "cancel")
	return s
}

type ShareSnapSpec struct {
	_JsonApiSpec[types.ShareSnapResult, protocol.ShareSnapResp]
}

func (s *ShareSnapSpec) Init(
	shareCode, receiveCode string, offset, limit int, dirId string,
) *ShareSnapSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/share/snap")
	s.query.Set("share_code", shareCode).
		Set("receive_code", receiveCode).
		Set("cid", dirId).
		SetInt("offset", offset).
		SetInt("limit", limit)
	return s
}

type ShareReceiveSpec struct {
	_VoidApiSpec
}

func (s *ShareReceiveSpec) Init(
	userId, shareCode, receiveCode string,
	fileIds []string, receiveDirId string,
) *ShareReceiveSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/share/receive")
	s.form.Set("user_id", userId).
		Set("share_code", shareCode).
		Set("receive_code", receiveCode).
		Set("file_id", strings.Join(fileIds, ","))
	if receiveDirId != "" {
		s.form.Set("cid", receiveDirId)
	}
	return s
}
