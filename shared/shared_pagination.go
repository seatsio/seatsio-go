package shared

type PaginationParams struct {
	PageSize    *int
	QueryParams map[string]string
}

type PaginationParamsOption func(Params *PaginationParams)

type paginationNS struct{}

var Pagination paginationNS

func (paginationNS) newParams() *PaginationParams {
	return &PaginationParams{QueryParams: map[string]string{}}
}

func (paginationNS) PageSize(pageSize int) PaginationParamsOption {
	return func(params *PaginationParams) {
		params.PageSize = &pageSize
	}
}

func (paginationNS) QueryParam(key string, value string) PaginationParamsOption {
	return func(params *PaginationParams) {
		params.QueryParams[key] = value
	}
}

func (paginationNS) QueryParams(queryParams map[string]string) PaginationParamsOption {
	return func(params *PaginationParams) {
		for key, value := range queryParams {
			params.QueryParams[key] = value
		}
	}
}
