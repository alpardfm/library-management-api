package query

import (
	"strings"

	"github.com/alpardfm/library-management-api/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type ListParams struct {
	Page   int
	Limit  int
	Search string
	Sort   string
}

type ListOptions struct {
	DefaultPage  int
	DefaultLimit int
	MaxLimit     int
	AllowedSorts map[string]string
	DefaultSort  string
}

func ParseListParams(c *gin.Context, opts ListOptions) (ListParams, error) {
	page, err := parsePositiveInt(c.Query("page"), opts.DefaultPage)
	if err != nil {
		return ListParams{}, apperror.BadRequest("invalid page query")
	}

	limit, err := parsePositiveInt(c.Query("limit"), opts.DefaultLimit)
	if err != nil {
		return ListParams{}, apperror.BadRequest("invalid limit query")
	}
	if opts.MaxLimit > 0 && limit > opts.MaxLimit {
		limit = opts.MaxLimit
	}

	search := strings.TrimSpace(c.Query("search"))

	sort := strings.TrimSpace(c.Query("sort"))
	if sort == "" {
		sort = opts.DefaultSort
	}
	if sort != "" {
		if _, ok := opts.AllowedSorts[sort]; !ok {
			return ListParams{}, apperror.BadRequest("invalid sort query")
		}
	}

	return ListParams{
		Page:   page,
		Limit:  limit,
		Search: search,
		Sort:   sort,
	}, nil
}
