package api

type UserInfoResult struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"face"`
	IsVip     int    `json:"vip"`
}

type UserInfoSpec struct {
	_StandardApiSpec[UserInfoResult]
}

func (s *UserInfoSpec) Init() *UserInfoSpec {
	s._StandardApiSpec.Init("https://my.115.com/?ct=ajax&ac=nav")
	return s
}
