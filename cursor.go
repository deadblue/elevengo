package elevengo

import "github.com/deadblue/elevengo/internal/webapi"

type FileCursor struct {
	// Transaction
	tx string
	// Cursor parameters
	offset int
	total  int
	// Sort parameters
	order string
	asc   int
}

func (c *FileCursor) checkTransaction(tx string) error {
	if c.tx != tx && c.tx != "" {
		return webapi.ErrInvalidCursor
	} else if c.tx == "" {
		// Initialize cursor
		c.tx = tx
		c.offset = 0
		c.order = "user_ptime"
		c.asc = 0
	}
	return nil
}

func (c *FileCursor) HasMore() bool {
	return c.tx == "" || c.offset < c.total
}

func (c *FileCursor) Total() int {
	return c.total
}

func (c *FileCursor) Remain() int {
	return c.total - c.offset
}
