package elevengo

type ClearFlag int

const (
	ClearComplete ClearFlag = iota
	ClearAll
	ClearFailed
	ClearRunning
	ClearCompleteAndDelete
	ClearAllAndDelete
)
