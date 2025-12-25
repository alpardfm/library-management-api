// internal/models/user.go
package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleLibrarian UserRole = "librarian"
	RoleMember    UserRole = "member"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(20);default:'member'" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	BorrowRecords []BorrowRecord `gorm:"foreignKey:UserID" json:"borrow_records,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}

	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	if u.PasswordHash == "" {
		return fmt.Errorf("password is required")
	}

	if strings.Contains(u.Email, "@") == false {
		return fmt.Errorf("invalid email format")
	}

	// Additional validation can be added here (e.g., email format)
	return nil
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}
