package webapi

type UserInfoData struct {
	UserId     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"face"`

	Device      int   `json:"device"`
	Rank        int   `json:"rank"`
	VipFlag     int   `json:"vip"`
	VipExpire   int64 `json:"expire"`
	VipForever  int   `json:"forever"`
	Global      int   `json:"global"`
	IsPrivilege bool  `json:"is_privilege"`
}
