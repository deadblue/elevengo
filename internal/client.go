package internal

type UserInfo struct {
	UserId string
}

type OfflineToken struct {
	Sign        string
	Time        int64
	QuotaTotal  int
	QuotaRemain int
}

type _BaseResult struct {
	State    bool   `json:"state"`
	ErrorNum int    `json:"errNo"`
	Error    string `json:"error"`
}
