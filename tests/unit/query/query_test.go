package query_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/pkg/query"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseListParams_DefaultsAndClamp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/books?limit=999&search=golang", nil)

	params, err := query.ParseListParams(c, query.ListOptions{
		DefaultPage:  1,
		DefaultLimit: 10,
		MaxLimit:     100,
		DefaultSort:  "created_at_desc",
		AllowedSorts: map[string]string{
			"created_at_desc": "created_at DESC",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 100, params.Limit)
	assert.Equal(t, "golang", params.Search)
	assert.Equal(t, "created_at_desc", params.Sort)
}

func TestParseListParams_InvalidQueries(t *testing.T) {
	tests := []string{
		"/books?page=abc",
		"/books?page=0",
		"/books?limit=zero",
		"/books?limit=0",
		"/books?sort=weird",
	}

	for _, url := range tests {
		t.Run(url, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", url, nil)

			_, err := query.ParseListParams(c, query.ListOptions{
				DefaultPage:  1,
				DefaultLimit: 10,
				MaxLimit:     100,
				DefaultSort:  "created_at_desc",
				AllowedSorts: map[string]string{
					"created_at_desc": "created_at DESC",
				},
			})

			assert.Error(t, err)
		})
	}
}

func TestTotalPages(t *testing.T) {
	assert.Equal(t, int64(0), query.TotalPages(0, 10))
	assert.Equal(t, int64(1), query.TotalPages(1, 10))
	assert.Equal(t, int64(3), query.TotalPages(21, 10))
}
