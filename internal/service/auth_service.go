// internal/service/auth_service.go
package service

import (
	"errors"
	"time"

	"library-management-api/internal/dto"
	"library-management-api/internal/models"
	"library-management-api/internal/repository"
	"library-management-api/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*models.User, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*models.User, error) {
	// Check if username exists
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	existingUser, _ = s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.UserRole(req.Role),
		IsActive:     true,
	}

	if user.Role == "" {
		user.Role = models.RoleMember
	}

	// Validate role
	switch user.Role {
	case models.RoleAdmin, models.RoleLibrarian, models.RoleMember:
		// Valid roles
	default:
		user.Role = models.RoleMember
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by username or email
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		// Try email
		user, err = s.userRepo.FindByEmail(req.Username)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dto.LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Role = string(user.Role)

	return response, nil
}

func (s *authService) GenerateToken(user *models.User) (string, error) {
	claims := auth.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString, s.jwtSecret)
}
