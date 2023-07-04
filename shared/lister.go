package shared

import "strconv"

type Lister[T interface{}] struct {
	PageFetcher *PageFetcher[T]
}

func (lister *Lister[T]) All(opts ...PaginationParamsOption) ([]T, error) {
	firstPage, err := lister.ListFirstPage(opts...)
	if err != nil {
		return nil, err
	}
	result := firstPage.Items
	currentPage := firstPage
	for currentPage.NextPageStartsAfter != 0 {
		currentPage, err = lister.ListPageAfter(currentPage.NextPageStartsAfter, opts...)
		if err != nil {
			return nil, err
		}
		result = append(result, currentPage.Items...)
	}
	return result, nil
}

func (lister *Lister[T]) ListFirstPage(opts ...PaginationParamsOption) (*Page[T], error) {
	return lister.PageFetcher.fetchPage(opts...)
}

func (lister *Lister[T]) ListPageAfter(id int64, opts ...PaginationParamsOption) (*Page[T], error) {
	newOpts := append(opts, Pagination.QueryParam("start_after_id", strconv.FormatInt(id, 10)))
	return lister.PageFetcher.fetchPage(newOpts...)
}

func (lister *Lister[T]) ListPageBefore(id int64, opts ...PaginationParamsOption) (*Page[T], error) {
	newOpts := append(opts, Pagination.QueryParam("end_before_id", strconv.FormatInt(id, 10)))
	return lister.PageFetcher.fetchPage(newOpts...)
}
