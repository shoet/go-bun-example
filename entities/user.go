package entities

import "github.com/uptrace/bun"

// UserはBookを持たないか、複数持つ
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID   int64 `bun:",pk,autoincrement"`
	Name string
	Book []Book `bun:"rel:has-many"`
}
