package webapi

type UserInfoData struct {
	UserId     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"face"`

	// Unused fields
	Device      int   `json:"device"`
	Rank        int   `json:"rank"`
	VipFlag     int   `json:"vip"`
	VipExpire   int64 `json:"expire"`
	VipForever  int   `json:"forever"`
	Global      int   `json:"global"`
	// IsPrivilege bool       `json:"is_privilege"`
	// Privilege   *Privilege `json:"privilege,omitempty"`
}

// type Privilege struct {
// 	State  bool `json:"state"`
// 	Start  int  `json:"start"`
// 	Expire int  `json:"expire"`
// }
