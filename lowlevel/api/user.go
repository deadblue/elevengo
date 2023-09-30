package api

import "github.com/deadblue/elevengo/internal/apibase"

type UserInfoResult struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"face"`
	IsVip     int    `json:"vip"`
}

type UserInfoSpec struct {
	apibase.JsonApiSpec[UserInfoResult, apibase.StandardResp]
}

func (s *UserInfoSpec) Init() *UserInfoSpec {
	s.JsonApiSpec.Init("https://my.115.com/?ct=ajax&ac=nav")
	return s
}
