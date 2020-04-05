package elevengo

// Cursor is a parameter using in list methods, such as
// * Agent.FileList
// * Agent.FileSearch
// * Agent.OfflineList
// These methods can not return the whole result in one time.
type Cursor interface {

	// Return true if there
	HasMore() bool

	// Move the cursor to fetch next page.
	Next()

	// Return the total count of the data.
	Total() int
}
