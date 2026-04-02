package query

import "strconv"

func parsePositiveInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if parsed < 1 {
		return 0, strconv.ErrSyntax
	}

	return parsed, nil
}

func TotalPages(total int64, limit int) int64 {
	if limit <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}

	pages := total / int64(limit)
	if total%int64(limit) != 0 {
		pages++
	}
	return pages
}
