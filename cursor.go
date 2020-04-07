package elevengo

// Cursor is a parameter using in list methods, such as
// * Agent.FileList
// * Agent.FileSearch
// * Agent.OfflineList
// These methods can not return the whole result in one time.
type Cursor interface {

	// Return true if the cursor has not been used or there is some data remain.
	HasMore() bool

	// Move cursor to the start of the remaining data.
	Next()

	// Return total amount of the data.
	Total() int
}
