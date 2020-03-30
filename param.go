package elevengo

const (
	pageSizeMin     = 3
	pageSizeMax     = 1000
	pageSizeDefault = 100
)

// Page parameter for `Client.FileList()` and `Client.FileSearch()`
type PageParam struct {
	index, size int
}

// The page size should be in [3, 1000]
func (p *PageParam) Size(size int) *PageParam {
	if size < pageSizeMin {
		p.size = pageSizeMin
	} else if size > pageSizeMax {
		p.size = pageSizeMax
	} else {
		p.size = size
	}
	return p
}

// Go to next page
func (p *PageParam) Next() *PageParam {
	p.index += 1
	return p
}

// Go to previous page
func (p *PageParam) Prev() *PageParam {
	if p.index > 0 {
		p.index -= 1
	}
	return p
}

// Go to specific page
func (p *PageParam) Goto(num int) *PageParam {
	if num > 0 {
		p.index = num - 1
	}
	return p
}
func (p *PageParam) limit() int {
	if p.size == 0 {
		p.size = pageSizeDefault
	}
	return p.size
}
func (p *PageParam) offset() int {
	return p.index * p.limit()
}

// Sort parameter for `Client.FileList()`
type SortParam struct {
	flag string
	asc  bool
}

func (p *SortParam) ByTime() *SortParam {
	p.flag = "user_ptime"
	return p
}
func (p *SortParam) ByName() *SortParam {
	p.flag = "file_name"
	return p
}
func (p *SortParam) BySize() *SortParam {
	p.flag = "file_size"
	return p
}
func (p *SortParam) Asc() *SortParam {
	p.asc = true
	return p
}
func (p *SortParam) Desc() *SortParam {
	p.asc = false
	return p
}
