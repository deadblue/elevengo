package types

type UserInfoResult struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"face"`
	IsVip     int    `json:"vip"`
}
