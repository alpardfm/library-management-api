package response

import (
	"errors"
	"net/http"

	"github.com/alpardfm/library-management-api/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Meta    any        `json:"meta,omitempty"`
}

func Success(c *gin.Context, status int, message string, data any, meta any) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, err error) {
	status, body := MapError(err)
	c.JSON(status, Envelope{
		Success: false,
		Message: body.Message,
		Error:   body,
	})
}

func MapError(err error) (int, *ErrorBody) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return statusForCode(appErr.Code), &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		}
	}

	return http.StatusInternalServerError, &ErrorBody{
		Code:    apperror.CodeInternal,
		Message: "internal server error",
	}
}

func statusForCode(code string) int {
	switch code {
	case apperror.CodeBadRequest:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
