// internal/handler/borrow_handler.go
package handler

import (
	"net/http"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/alpardfm/library-management-api/pkg/apperror"
	"github.com/alpardfm/library-management-api/pkg/query"
	httpresponse "github.com/alpardfm/library-management-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	borrowService service.BorrowService
}

func NewBorrowHandler(borrowService service.BorrowService) *BorrowHandler {
	return &BorrowHandler{borrowService: borrowService}
}

func (h *BorrowHandler) BorrowBook(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	borrowRecord, err := h.borrowService.BorrowBook(userID, req)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusCreated, "Book borrowed successfully", borrowRecord, nil)
}

func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var req dto.ReturnBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	borrowRecord, fine, err := h.borrowService.ReturnBook(userID, role, req)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	data := gin.H{
		"borrow_record": borrowRecord,
	}
	if fine > 0 {
		data["fine"] = fine
	}

	httpresponse.Success(c, http.StatusOK, "Book returned successfully", data, nil)
}

func (h *BorrowHandler) GetMyBorrows(c *gin.Context) {
	userID := c.GetUint("user_id")
	params, err := query.ParseListParams(c, query.ListOptions{
		DefaultPage:  1,
		DefaultLimit: 10,
		MaxLimit:     50,
		DefaultSort:  "created_at_desc",
		AllowedSorts: map[string]string{
			"created_at_desc": "created_at DESC",
			"created_at_asc":  "created_at ASC",
			"due_date_asc":    "due_date ASC",
			"due_date_desc":   "due_date DESC",
		},
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	borrows, total, err := h.borrowService.GetUserBorrows(userID, params.Page, params.Limit, params.Sort)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "", borrows, gin.H{
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": query.TotalPages(total, params.Limit),
		"sort":        params.Sort,
	})
}

func (h *BorrowHandler) GetActiveBorrows(c *gin.Context) {
	params, err := query.ParseListParams(c, query.ListOptions{
		DefaultPage:  1,
		DefaultLimit: 10,
		MaxLimit:     50,
		DefaultSort:  "due_date_asc",
		AllowedSorts: map[string]string{
			"created_at_desc": "created_at DESC",
			"created_at_asc":  "created_at ASC",
			"due_date_asc":    "due_date ASC",
			"due_date_desc":   "due_date DESC",
		},
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	borrows, total, err := h.borrowService.GetActiveBorrows(params.Page, params.Limit, params.Sort)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "", borrows, gin.H{
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": query.TotalPages(total, params.Limit),
		"sort":        params.Sort,
	})
}

func (h *BorrowHandler) GetOverdueBorrows(c *gin.Context) {
	params, err := query.ParseListParams(c, query.ListOptions{
		DefaultPage:  1,
		DefaultLimit: 10,
		MaxLimit:     50,
		DefaultSort:  "due_date_asc",
		AllowedSorts: map[string]string{
			"created_at_desc": "created_at DESC",
			"created_at_asc":  "created_at ASC",
			"due_date_asc":    "due_date ASC",
			"due_date_desc":   "due_date DESC",
		},
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	borrows, total, err := h.borrowService.GetOverdueBorrows(params.Page, params.Limit, params.Sort)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "", borrows, gin.H{
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": query.TotalPages(total, params.Limit),
		"sort":        params.Sort,
	})
}
