package entities

import "github.com/uptrace/bun"

// Bookはユーザーを一つあるいは、持たない
type Book struct {
	bun.BaseModel `bun:"table:books"`

	ID    int64 `bun:",pk,autoincrement"`
	Title string
	User  User `bun:"rel:has-one"`
}
