// internal/models/book.go
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ISBN            string    `gorm:"uniqueIndex;size:13;not null" json:"isbn"`
	Title           string    `gorm:"size:255;not null" json:"title"`
	Author          string    `gorm:"size:255;not null" json:"author"`
	Publisher       string    `gorm:"size:100" json:"publisher,omitempty"`
	PublicationYear int       `json:"publication_year,omitempty"`
	Genre           string    `gorm:"size:50" json:"genre,omitempty"`
	Description     string    `gorm:"type:text" json:"description,omitempty"`
	TotalCopies     int       `gorm:"default:1" json:"total_copies"`
	AvailableCopies int       `gorm:"default:1" json:"available_copies"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	BorrowRecords []BorrowRecord `gorm:"foreignKey:BookID" json:"borrow_records,omitempty"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	// Initialize available copies
	if b.AvailableCopies == 0 && b.TotalCopies > 0 {
		b.AvailableCopies = b.TotalCopies
	}
	return nil
}

func (b *Book) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

// CanBorrow checks if book is available for borrowing
func (b *Book) CanBorrow() bool {
	return b.AvailableCopies > 0
}

// Borrow reduces available copies
func (b *Book) Borrow() error {
	if !b.CanBorrow() {
		return fmt.Errorf("book not available for borrowing")
	}
	b.AvailableCopies--
	return nil
}

// Return increases available copies
func (b *Book) Return() {
	if b.AvailableCopies < b.TotalCopies {
		b.AvailableCopies++
	}
}
