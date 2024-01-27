package entities

type Author struct {
	ID   int64  `bun:"id,pk,autoincrement"`
	Name string `bun:"name,unique"`
}
