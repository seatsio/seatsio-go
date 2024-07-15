package shared

import (
	"github.com/imroc/req/v3"
	"strconv"
)

type PageFetcher[T interface{}] struct {
	Client      *req.Client
	Url         string
	UrlParams   map[string]string
	QueryParams map[string]string
}

type Page[T interface{}] struct {
	Items                  []T   `json:"items"`
	NextPageStartsAfter    int64 `json:"next_page_starts_after"`
	PreviousPageEndsBefore int64 `json:"previous_page_ends_before"`
}

type PageJson[T interface{}] struct {
	Items                  []T    `json:"items"`
	NextPageStartsAfter    string `json:"next_page_starts_after"`
	PreviousPageEndsBefore string `json:"previous_page_ends_before"`
}

func (pageFetcher *PageFetcher[T]) fetchPage(opts ...PaginationParamsOption) (*Page[T], error) {
	paginationParams := Pagination.newParams()
	for _, opt := range opts {
		opt(paginationParams)
	}
	var page PageJson[T]
	request := pageFetcher.Client.R().
		SetSuccessResult(&page)

	if paginationParams.PageSize != nil {
		request.SetQueryParam("limit", strconv.Itoa(*paginationParams.PageSize))
	}
	if paginationParams.QueryParams != nil {
		for key, value := range paginationParams.QueryParams {
			request.AddQueryParam(key, value)
		}
	}

	if paginationParams.QueryParamsArrays != nil {
		for key, value := range paginationParams.QueryParamsArrays {
			request.AddQueryParams(key, value...)
		}
	}

	if pageFetcher.QueryParams != nil {
		for key, value := range pageFetcher.QueryParams {
			request.AddQueryParam(key, value)
		}
	}

	for key, value := range pageFetcher.UrlParams {
		request.SetPathParam(key, value)
	}

	result, err := request.Get(pageFetcher.Url)
	_, err = AssertOk(result, err, &page)
	if err != nil {
		return nil, err
	}

	nextPageStartsAfterInt, err := optionalIdToInt(page.NextPageStartsAfter)
	if err != nil {
		return nil, err
	}

	previousPageEndsBeforeInt, err := optionalIdToInt(page.PreviousPageEndsBefore)
	if err != nil {
		return nil, err
	}

	return &Page[T]{page.Items, nextPageStartsAfterInt, previousPageEndsBeforeInt}, nil
}

func optionalIdToInt(id string) (int64, error) {
	if id == "" {
		return 0, nil
	}
	return strconv.ParseInt(id, 10, 64)
}

func ToSort(sortField string, sortDirection string) string {
	if sortField == "" {
		return ""
	}
	if sortDirection == "" {
		return sortField
	}
	return sortField + ":" + sortDirection
}
