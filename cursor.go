package elevengo

/*
Due to the upstream API restriction, some list methods can not get the whole data
in one request. Cursor is design for these methods, to get data page-by-page.

There are two types of cursors for different methods:

* File cursor: created by FileCursor(), used in Agent.FileList() and Agent.FileSearch().

* Offline cursor: created by OfflineCursor(), used in Agent.OfflineList().

A typical usage (call FileList for example):

	// Assume "agent" is an Agent instance, and "parentId" is the directory ID
	// where you want to get file list from.
	for cursor := FileCursor(); cursor.HasMore(); cursor.Next() {
		if files, err := agent.FileList(parentId, cursor); err != nil {
			// handle error
			break
		} else {
			// deal with the files
		}
	}
*/
type Cursor interface {

	// Return true if the cursor has not been used or there is some data remain.
	HasMore() bool

	// Move cursor to the start of the remaining data.
	Next()

	// Return total amount of the data.
	Total() int
}
