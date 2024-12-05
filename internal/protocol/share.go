package protocol

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type ShareListResp struct {
	BasicResp

	Count int             `json:"count"`
	List  json.RawMessage `json:"list"`
}

func (r *ShareListResp) Extract(v *types.ShareListResult) (err error) {
	if err = json.Unmarshal(r.List, &v.Items); err != nil {
		return
	}
	v.Count = r.Count
	return
}

type ShareFileProto struct {
	FileId     string         `json:"fid"`
	DirId      util.IntNumber `json:"cid"`
	ParentId   string         `json:"pid"`
	Name       string         `json:"n"`
	IsFile     int            `json:"fc"`
	Size       int64          `json:"s"`
	Sha1       string         `json:"sha"`
	CreateTime util.IntNumber `json:"t"`

	IsVideo         int `json:"iv"`
	VideoDefinition int `json:"vdi"`
	AudioPlayLong   int `json:"audio_play_long"`
	VideoPlayLong   int `json:"play_long"`
}

type ShareSnapProto struct {
	UserInfo struct {
		UserId    util.IntNumber `json:"user_id"`
		UserName  string         `json:"user_name"`
		AvatarUrl string         `json:"face"`
	} `json:"userinfo"`
	UserAppeal struct {
		CanAppeal       int `json:"can_appeal"`
		CanShareAppeal  int `json:"can_share_appeal"`
		CanGlobalAppeal int `json:"can_global_appeal"`
		PopupAppealPage int `json:"popup_appeal_page"`
	} `json:"user_appeal"`
	ShareInfo struct {
		SnapId              string         `json:"snap_id"`
		ShareTitle          string         `json:"share_title"`
		ShareState          string         `json:"share_state"`
		ShareSize           util.IntNumber `json:"file_size"`
		ReceiveCode         string         `json:"receive_code"`
		ReceiveCount        util.IntNumber `json:"receive_count"`
		CreateTime          util.IntNumber `json:"create_time"`
		ExpireTime          util.IntNumber `json:"expire_time"`
		AutoRenewal         string         `json:"auto_renewal"`
		AutoFillReceiveCode string         `json:"auto_fill_recvcode"`
		CanReport           int            `json:"can_report"`
		CanNotice           int            `json:"can_notice"`
		HaveVioFile         int            `json:"have_vio_file"`
	} `json:"shareinfo"`
	ShareState string `json:"share_state"`

	Count int               `json:"count"`
	Files []*ShareFileProto `json:"list"`
}

type ShareSnapResp struct {
	BasicResp

	Data ShareSnapProto `json:"data"`
}

func (r *ShareSnapResp) Extract(v *types.ShareSnapResult) (err error) {
	data := r.Data
	v.SnapId = data.ShareInfo.SnapId
	v.UserId = data.UserInfo.UserId.Int()
	v.ShareTitle = data.ShareInfo.ShareTitle
	v.ShareState = data.ShareInfo.ShareSize.Int()
	v.ReceiveCount = data.ShareInfo.ReceiveCount.Int()
	v.CreateTime = data.ShareInfo.CreateTime.Int64()
	v.ExpireTime = data.ShareInfo.ExpireTime.Int64()

	v.TotalSize = data.ShareInfo.ShareSize.Int64()
	v.FileCount = data.Count
	v.Files = make([]*types.ShareFileInfo, len(data.Files))
	for i, f := range data.Files {
		fileId := ""
		if f.IsFile == 1 {
			fileId = f.FileId
		} else {
			fileId = f.DirId.String()
		}
		v.Files[i] = &types.ShareFileInfo{
			FileId:     fileId,
			IsDir:      f.IsFile == 0,
			Name:       f.Name,
			Size:       f.Size,
			Sha1:       f.Sha1,
			CreateTime: f.CreateTime.Int64(),

			IsVideo:         f.IsVideo != 0,
			VideoDefinition: f.VideoDefinition,
			MediaDuration:   max(f.AudioPlayLong, f.VideoPlayLong),
		}
	}
	return
}
