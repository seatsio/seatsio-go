package shared

type paginationParams struct {
	PageSize    *int
	QueryParams map[string]string
}

type PaginationParamsOption func(Params *paginationParams)

type paginationNS struct{}

var Pagination paginationNS

func (paginationNS) newParams() *paginationParams {
	return &paginationParams{QueryParams: map[string]string{}}
}

func (paginationNS) PageSize(pageSize int) PaginationParamsOption {
	return func(params *paginationParams) {
		params.PageSize = &pageSize
	}
}

func (paginationNS) QueryParam(key string, value string) PaginationParamsOption {
	return func(params *paginationParams) {
		params.QueryParams[key] = value
	}
}

func (paginationNS) QueryParams(queryParams map[string]string) PaginationParamsOption {
	return func(params *paginationParams) {
		for key, value := range queryParams {
			params.QueryParams[key] = value
		}
	}
}
