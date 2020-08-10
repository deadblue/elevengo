package elevengo

/*
Due to the upstream API restriction, some list methods can not get the whole data
in one request. Cursor is design for these methods, to get data page-by-page.
*/
type Cursor interface {

	// Return true if the cursor has not been used or there is some data remain.
	HasMore() bool

	// Move cursor to the start of the remaining data.
	Next()

	// Return total amount of the data.
	Total() int
}

/*
File cursor for FileList and FileSearch, developer should
create it through "FileCursor()".
*/
type fileCursor struct {
	used   bool
	order  string
	asc    int
	offset int
	limit  int
	total  int
}

func (c *fileCursor) HasMore() bool {
	return !c.used || c.offset < c.total
}
func (c *fileCursor) Next() {
	c.offset += c.limit
}
func (c *fileCursor) Total() int {
	return c.total
}

// Create a cursor for "Agent.FileList()" and "Agent.FileSearch()".
func FileCursor() Cursor {
	return &fileCursor{
		used:   false,
		order:  "user_ptime",
		asc:    0,
		offset: 0,
		limit:  fileDefaultLimit,
		total:  0,
	}
}

/*
File cursor for OfflineList, developer should create it through "OfflineCursor()".
*/
type offlineCursor struct {
	used      bool
	page      int
	pageCount int
	total     int
}

func (c *offlineCursor) HasMore() bool {
	return !c.used || c.page < c.pageCount
}
func (c *offlineCursor) Next() {
	c.page += 1
}
func (c *offlineCursor) Total() int {
	return c.total
}

// Create a cursor for "Agent.OfflineList()".
func OfflineCursor() Cursor {
	return &offlineCursor{
		used:      false,
		page:      1,
		pageCount: 0,
		total:     0,
	}
}
