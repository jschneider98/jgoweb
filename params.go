package jgoweb

type SearchParams struct {
	Query string
	Limit uint64
}

//
func NewSearchParams() *SearchParams {
	return &SearchParams{}
}
