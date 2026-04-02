// internal/handler/book_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/alpardfm/library-management-api/pkg/apperror"
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	books, total, err := h.bookService.ListBooks(page, limit, search)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "", books, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
	})
}
