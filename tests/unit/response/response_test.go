package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/alpardfm/library-management-api/pkg/apperror"
	"github.com/alpardfm/library-management-api/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestMapError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedType string
	}{
		{name: "bad request", err: apperror.BadRequest("bad input"), expectedCode: http.StatusBadRequest, expectedType: apperror.CodeBadRequest},
		{name: "unauthorized", err: apperror.Unauthorized("invalid credentials"), expectedCode: http.StatusUnauthorized, expectedType: apperror.CodeUnauthorized},
		{name: "forbidden", err: apperror.Forbidden("forbidden"), expectedCode: http.StatusForbidden, expectedType: apperror.CodeForbidden},
		{name: "not found", err: apperror.NotFound("book"), expectedCode: http.StatusNotFound, expectedType: apperror.CodeNotFound},
		{name: "conflict", err: apperror.Conflict("already exists"), expectedCode: http.StatusConflict, expectedType: apperror.CodeConflict},
		{name: "unknown", err: errors.New("boom"), expectedCode: http.StatusInternalServerError, expectedType: apperror.CodeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := response.MapError(tt.err)
			assert.Equal(t, tt.expectedCode, status)
			assert.Equal(t, tt.expectedType, body.Code)
		})
	}
}
