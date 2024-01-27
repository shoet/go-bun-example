package entities

// Bookはユーザーを一つあるいは、持たない
type Book struct {
	ID       int64   `bun:"id,pk,autoincrement"`
	Title    string  `bun:"title"`
	AuthorID int64   `bun:"author_id"`
	Author   *Author `bun:"rel:belongs-to,join:author_id=id"`
}
