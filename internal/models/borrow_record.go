// internal/models/borrow_record.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type BorrowStatus string

const (
	StatusBorrowed BorrowStatus = "borrowed"
	StatusReturned BorrowStatus = "returned"
	StatusOverdue  BorrowStatus = "overdue"
)

type BorrowRecord struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	UserID     uint         `gorm:"not null" json:"user_id"`
	BookID     uint         `gorm:"not null" json:"book_id"`
	BorrowDate time.Time    `gorm:"not null" json:"borrow_date"`
	DueDate    time.Time    `gorm:"not null" json:"due_date"`
	ReturnDate *time.Time   `json:"return_date,omitempty"`
	Status     BorrowStatus `gorm:"type:varchar(20);default:'borrowed'" json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

func (br *BorrowRecord) BeforeCreate(tx *gorm.DB) error {
	br.CreatedAt = time.Now()
	br.UpdatedAt = time.Now()

	// Set borrow date to now if not set
	if br.BorrowDate.IsZero() {
		br.BorrowDate = time.Now()
	}

	// Set due date to 14 days from borrow date if not set
	if br.DueDate.IsZero() {
		br.DueDate = br.BorrowDate.Add(14 * 24 * time.Hour) // 14 days
	}

	// Set initial status
	if br.Status == "" {
		br.Status = StatusBorrowed
	}

	return nil
}

func (br *BorrowRecord) BeforeUpdate(tx *gorm.DB) error {
	br.UpdatedAt = time.Now()

	// Update status based on return date
	if br.ReturnDate != nil && br.Status != StatusReturned {
		br.Status = StatusReturned
	} else if br.ReturnDate == nil && time.Now().After(br.DueDate) {
		br.Status = StatusOverdue
	}

	return nil
}

// IsOverdue checks if the borrow is overdue
func (br *BorrowRecord) IsOverdue() bool {
	return br.Status == StatusOverdue || (br.ReturnDate == nil && time.Now().After(br.DueDate))
}

// CalculateFine calculates fine for overdue books (Rp 1000 per day)
func (br *BorrowRecord) CalculateFine() int {
	if br.ReturnDate != nil || !br.IsOverdue() {
		return 0
	}

	overdueDays := int(time.Since(br.DueDate).Hours() / 24)
	if overdueDays < 0 {
		return 0
	}

	return overdueDays * 1000 // Rp 1000 per day
}
