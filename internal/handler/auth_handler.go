// internal/handler/auth_handler.go
package handler

import (
	"net/http"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/alpardfm/library-management-api/pkg/apperror"
	httpresponse "github.com/alpardfm/library-management-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.Error(c, apperror.BadRequest(err.Error()))
		return
	}

	loginResponse, err := h.authService.Login(req)
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.Success(c, http.StatusOK, "Login successful", loginResponse, nil)
}
