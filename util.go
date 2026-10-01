package dtrack

import "fmt"

// FetchAll is a convenience function to retrieve all items of a paginated API resource.
func FetchAll[T any](pageFetchFunc func(po PageOptions) (Page[T], error)) (items []T, err error) {
	err = ForEach(pageFetchFunc, func(item T) error {
		items = append(items, item)
		return nil
	})

	return
}

// ForEach is a convenience function to perform an action on every item of a paginated API resource.
func ForEach[T any](pageFetchFunc func(po PageOptions) (Page[T], error), handlerFunc func(item T) error) (err error) {
	const pageSize = 50

	var (
		page       Page[T]
		pageNumber = 1
		itemsSeen  = 0
	)

	for {
		page, err = pageFetchFunc(PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			break
		}

		for i := range page.Items {
			err = handlerFunc(page.Items[i])
			if err != nil {
				return fmt.Errorf("failed to handle item %d on page %d: %w", i+1, pageNumber, err)
			}
		}

		itemsSeen += len(page.Items)
		if len(page.Items) == 0 || itemsSeen >= page.TotalCount {
			break
		}

		pageNumber++
	}

	return
}

func OptionalBoolOf(value bool) *bool {
	return &value
}

func OptionalBool() *bool {
	return nil
}

// FetchAllV5 is a convenience function to retrieve all items of a token-paginated API v2 resource.
func FetchAllV5[T any](pageFetchFunc func(pageToken string, limit uint32) (PageV5[T], error)) (items []T, err error) {
	err = ForEachV5(pageFetchFunc, func(item T) error {
		items = append(items, item)
		return nil
	})

	return
}

// ForEachV5 is a convenience function to perform an action on every item of a token-paginated API v2 resource.
func ForEachV5[T any](pageFetchFunc func(pageToken string, limit uint32) (PageV5[T], error), handlerFunc func(item T) error) (err error) {
	const limit = 100

	var (
		page      PageV5[T]
		pageToken = ""
	)

	for {
		page, err = pageFetchFunc(pageToken, limit)
		if err != nil {
			break
		}

		for i := range page.Items {
			err = handlerFunc(page.Items[i])
			if err != nil {
				return fmt.Errorf("failed to handle item %d with page token %q: %w", i+1, pageToken, err)
			}
		}

		if len(page.Items) == 0 || page.NextPageToken == "" {
			break
		}

		pageToken = page.NextPageToken
	}

	return
}
