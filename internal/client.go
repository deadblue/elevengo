package internal

type UserInfo struct {
	UserId   int
	UserName string
}

type OfflineToken struct {
	Sign string
	Time int64
}
