package internal

type UserInfo struct {
	UserId   string
	UserName string
}

type OfflineToken struct {
	Sign string
	Time int64
}
