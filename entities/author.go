package entities

type Author struct {
	ID   int64  `bun:"id,pk"`
	Name string `bun:"name"`
}
