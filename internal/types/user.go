package types

type UserInfoResult struct {
	BaseApiResult
	Data struct {
		UserId      int    `json:"user_id"`
		UserName    string `json:"user_name"`
		UserAvatar  string `json:"face"`
		Device      int    `json:"device"`
		Rank        int    `json:"rank"`
		VipFlag     int    `json:"vip"`
		VipExpire   int64  `json:"expire"`
		Global      int    `json:"global"`
		Forever     int    `json:"forever"`
		IsPrivilege bool   `json:"is_privilege"`
		Privilege   struct {
			State  bool `json:"state"`
			Start  int  `json:"start"`
			Expire int  `json:"expire"`
		} `json:"privilege"`
	} `json:"data"`
}
