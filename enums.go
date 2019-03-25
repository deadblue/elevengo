package elevengo

type OrderFlag string

const (
	OrderByTime OrderFlag = "user_ptime"
	OrderByName OrderFlag = "file_name"
	OrderBySize OrderFlag = "file_size"
)

type ClearFlag int

const (
	ClearComplete ClearFlag = iota
	ClearAll
	ClearFailed
	ClearRunning
	ClearCompleteAndDelete
	ClearAllAndDelete
)
