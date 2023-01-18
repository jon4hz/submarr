package httpclient

type Request struct {
	params        map[string]string
	page          int
	pageSize      int
	sortKey       string
	sortDirection SortDirection
}

type SortDirection string

const (
	Default    SortDirection = "default"
	Ascending  SortDirection = "ascending"
	Descending SortDirection = "descending"
)

func newRequest(opts ...RequestOpts) *Request {
	r := &Request{}
	for _, o := range opts {
		o(r)
	}
	return r
}

type RequestOpts func(*Request)

func WithParams(params map[string]string) RequestOpts {
	return func(r *Request) {
		r.params = params
	}
}

func WithPage(page int) RequestOpts {
	return func(r *Request) {
		r.page = page
	}
}

func WithPageSize(pageSize int) RequestOpts {
	return func(r *Request) {
		r.pageSize = pageSize
	}
}

func WithSortKey(sortKey string) RequestOpts {
	return func(r *Request) {
		r.sortKey = sortKey
	}
}

func WithSortDirection(sortDirection SortDirection) RequestOpts {
	return func(r *Request) {
		r.sortDirection = sortDirection
	}
}
