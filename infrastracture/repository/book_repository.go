package repository

import (
	"context"
	"fmt"

	"github.com/shoet/go-bun-example/entities"
	"github.com/uptrace/bun"
)

type BookRepository struct {
}

func NewBookRepository() (*BookRepository, error) {
	return &BookRepository{}, nil
}

// upsert relation
func (b *BookRepository) CreateBook(
	ctx context.Context, tx *bun.Tx, book *entities.Book,
) (bookID int64, err error) {
	author := book.Author
	if _, err := tx.
		NewInsert().
		Model(author).
		On("DUPLICATE KEY UPDATE").
		Returning("*").
		Exec(ctx); err != nil {
		return 0, fmt.Errorf("failed to create author: %w", err)
	}
	if err := tx.NewSelect().Model(author).Where("name = ?", author.Name).Scan(ctx); err != nil {
		return 0, fmt.Errorf("failed to get author: %w", err)
	}

	book.AuthorID = author.ID
	if _, err := tx.
		NewInsert().
		Model(book).
		Exec(ctx); err != nil {
		return 0, fmt.Errorf("failed to create book: %w", err)
	}
	return book.ID, nil
}
