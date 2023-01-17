package httpclient

type request struct {
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

func newRequest(opts ...RequestOpts) *request {
	r := &request{}
	for _, o := range opts {
		o(r)
	}
	return r
}

type RequestOpts func(*request)

func WithParams(params map[string]string) RequestOpts {
	return func(r *request) {
		r.params = params
	}
}

func WithPage(page int) RequestOpts {
	return func(r *request) {
		r.page = page
	}
}

func WithPageSize(pageSize int) RequestOpts {
	return func(r *request) {
		r.pageSize = pageSize
	}
}

func WithSortKey(sortKey string) RequestOpts {
	return func(r *request) {
		r.sortKey = sortKey
	}
}

func WithSortDirection(sortDirection SortDirection) RequestOpts {
	return func(r *request) {
		r.sortDirection = sortDirection
	}
}
