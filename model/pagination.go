package model

type PageItem struct {
	Title   string
	PageNum int
	Active  bool
	Url     string
}

type Pagination struct {
	Page     int
	PageSize int
	Data     interface{}
	Total    int64
	Items    []PageItem
}
