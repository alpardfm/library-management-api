// internal/dto/book.go
package dto

type CreateBookRequest struct {
	ISBN            string `json:"isbn" binding:"required,min=10,max=13"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	Publisher       string `json:"publisher,omitempty"`
	PublicationYear int    `json:"publication_year,omitempty" binding:"gte=1000,lte=2024"`
	Genre           string `json:"genre,omitempty"`
	Description     string `json:"description,omitempty"`
	TotalCopies     int    `json:"total_copies" binding:"gte=1"`
}

type UpdateBookRequest struct {
	Title           string `json:"title,omitempty"`
	Author          string `json:"author,omitempty"`
	Publisher       string `json:"publisher,omitempty"`
	PublicationYear int    `json:"publication_year,omitempty" binding:"omitempty,gte=1000,lte=2024"`
	Genre           string `json:"genre,omitempty"`
	Description     string `json:"description,omitempty"`
	TotalCopies     int    `json:"total_copies,omitempty" binding:"omitempty,gte=1"`
}

type BookResponse struct {
	ID              uint   `json:"id"`
	ISBN            string `json:"isbn"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	Publisher       string `json:"publisher,omitempty"`
	PublicationYear int    `json:"publication_year,omitempty"`
	Genre           string `json:"genre,omitempty"`
	Description     string `json:"description,omitempty"`
	TotalCopies     int    `json:"total_copies"`
	AvailableCopies int    `json:"available_copies"`
}
