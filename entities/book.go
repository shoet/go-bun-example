package entities

import "github.com/uptrace/bun"

// Bookはユーザーを一つあるいは、持たない
type Book struct {
	bun.BaseModel `bun:"table:books"`

	ID       int64 `bun:",pk,autoincrement"`
	Title    string
	AuthorID int64  `bun:"author_id"`
	Author   Author `bun:"rel:has-one,join:author_id=id"`
}
