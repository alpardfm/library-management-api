// internal/handler/book_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/alpardfm/library-management-api/pkg/apperror"
	"github.com/alpardfm/library-management-api/pkg/query"
	httpresponse "github.com/alpardfm/library-management-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	book, err := h.bookService.CreateBook(req)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusCreated, "Book created successfully", book, nil)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpresponse.Error(c, apperror.BadRequest("invalid book ID"))
		return
	}

	book, err := h.bookService.GetBookByID(uint(id))
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "", book, nil)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpresponse.Error(c, apperror.BadRequest("invalid book ID"))
		return
	}

	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	book, err := h.bookService.UpdateBook(uint(id), req)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "Book updated successfully", book, nil)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpresponse.Error(c, apperror.BadRequest("invalid book ID"))
		return
	}

	if err := h.bookService.DeleteBook(uint(id)); err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "Book deleted successfully", nil, nil)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	params, err := query.ParseListParams(c, query.ListOptions{
		DefaultPage:  1,
		DefaultLimit: 10,
		MaxLimit:     100,
		DefaultSort:  "created_at_desc",
		AllowedSorts: map[string]string{
			"created_at_desc": "created_at DESC",
			"created_at_asc":  "created_at ASC",
			"title_asc":       "title ASC",
			"title_desc":      "title DESC",
		},
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	books, total, err := h.bookService.ListBooks(params.Page, params.Limit, params.Search, params.Sort)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "", books, gin.H{
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": query.TotalPages(total, params.Limit),
		"sort":        params.Sort,
		"search":      params.Search,
	})
}
