// internal/dto/borrow.go
package dto

import "time"

type BorrowBookRequest struct {
	BookID  uint      `json:"book_id" binding:"required"`
	UserID  uint      `json:"user_id,omitempty"` // Admin bisa specify user lain
	DueDate time.Time `json:"due_date,omitempty"`
}

type ReturnBookRequest struct {
	BorrowRecordID uint `json:"borrow_record_id" binding:"required"`
}

type BorrowRecordResponse struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	BookID     uint       `json:"book_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	Status     string     `json:"status"`
	Fine       int        `json:"fine,omitempty"`

	// Nested objects
	User struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user,omitempty"`

	Book struct {
		ID     uint   `json:"id"`
		ISBN   string `json:"isbn"`
		Title  string `json:"title"`
		Author string `json:"author"`
	} `json:"book,omitempty"`
}
