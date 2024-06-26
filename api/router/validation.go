package router

import "voltaserve/service"

func IsValidSortBy(value string) bool {
	return value == "" ||
		value == service.SortByName ||
		value == service.SortByKind ||
		value == service.SortBySize ||
		value == service.SortByEmail ||
		value == service.SortByFullName ||
		value == service.SortByVersion ||
		value == service.SortByFrequency ||
		value == service.SortByDateCreated ||
		value == service.SortByDateModified
}

func IsValidSortOrder(value string) bool {
	return value == "" || value == service.SortOrderAsc || value == service.SortOrderDesc
}
