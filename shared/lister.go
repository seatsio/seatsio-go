package shared

import "strconv"

type Lister[T interface{}] struct {
	PageFetcher *PageFetcher[T]
}

func (lister *Lister[T]) All(pageSize int) ([]T, error) {
	firstPage, err := lister.ListFirstPage(pageSize)
	if err != nil {
		return nil, err
	}
	result := firstPage.Items
	currentPage := firstPage
	for currentPage.NextPageStartsAfter != 0 {
		currentPage, err = lister.ListPageAfter(currentPage.NextPageStartsAfter, pageSize)
		if err != nil {
			return nil, err
		}
		result = append(result, currentPage.Items...)
	}
	return result, nil
}

func (lister *Lister[T]) ListFirstPage(pageSize int) (*Page[T], error) {
	return lister.PageFetcher.fetchPage(pageSize, map[string]string{})
}

func (lister *Lister[T]) ListPageAfter(id int64, pageSize int) (*Page[T], error) {
	return lister.PageFetcher.fetchPage(pageSize, map[string]string{"start_after_id": strconv.FormatInt(id, 10)})
}

func (lister *Lister[T]) ListPageBefore(id int64, pageSize int) (*Page[T], error) {
	return lister.PageFetcher.fetchPage(pageSize, map[string]string{"end_before_id": strconv.FormatInt(id, 10)})
}
