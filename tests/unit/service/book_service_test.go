package service_test

import (
	"testing"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookService_UpdateBook_RejectsTotalCopiesBelowBorrowedCopies(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		ISBN:            "9781234567897",
		Title:           "Distributed Systems",
		Author:          "Author",
		TotalCopies:     5,
		AvailableCopies: 1,
	}

	req := dto.UpdateBookRequest{
		TotalCopies: 3,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()

	book, err := bookService.UpdateBook(1, req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "total copies cannot be less than borrowed copies", err.Error())
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestBookService_UpdateBook_RejectsInconsistentExistingStock(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		ISBN:            "9781234567897",
		Title:           "Distributed Systems",
		Author:          "Author",
		TotalCopies:     2,
		AvailableCopies: 3,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()

	book, err := bookService.UpdateBook(1, dto.UpdateBookRequest{Title: "New Title"})

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "book stock is inconsistent", err.Error())
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}
